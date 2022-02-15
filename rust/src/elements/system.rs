use std::{time::Duration, sync::{Arc}};
use bson::oid::ObjectId;
use llml::vec::EucVecd2;
use serde::{Serialize, Deserialize};
use turing_proc::Maybee;
use crate::{Star, Planet, cache::MongoDoc};
use std::hash::Hash;
use tokio::{sync::Mutex};

#[derive(Clone, Debug, Serialize, Deserialize, Maybee)]
pub struct PlanetSystem {
    #[serde(rename = "_id")]
    pub id: ObjectId,
    pub star: Star,
    pub planets: Vec<Planet>
}

impl PlanetSystem {
    pub fn new (star: Star, planets: Vec<Planet>) -> Self {
        PlanetSystem { id: ObjectId::new(), star, planets }
    }

    pub async fn simulate (&mut self, dt: Duration) {
        // Iterator over all possible planet pairs. This way we don't repeat calculations.
        let len = self.planets.len();
        let mut pairs = Vec::with_capacity((len * len - len) / 2);

        for i in 0..(len-1) {
            for j in i..len {
                pairs.push((&self.planets[i], &self.planets[j]));
            }
        }

        // Initialize planet acceleration by calculating its acceleration to the start. In this model, stars are immovable (always at coordinate origin)
        let mut interplanet_acc = Vec::<EucVecd2>::with_capacity(self.planets.len());
        for planet in self.planets.iter() {
            interplanet_acc.insert(planet.id, planet.calc_acc_star(&self.star));
        }

        // Calculate interplanet acceleration for each planet pair. Each calculation is done in a diferent thread
        let interplanet_acc = Arc::new(Mutex::new(interplanet_acc));
        let handles = pairs.into_iter().map(|(x, y)| {
            let acc_clone = interplanet_acc.clone();
            async move {
                let (acc_x, acc_y) = x.calc_acc(y);
                let mut lock = acc_clone.lock().await;
                lock[x.id] += acc_x;
                lock[y.id] += acc_y;
            }
        });

        // Process features concurrently
        futures::future::join_all(handles).await;

        // Apply changes
        let planets = self.planets.iter_mut();
        let lock = interplanet_acc.lock().await;
        planets.for_each(|planet| { planet.accelerate_and_travel(lock[planet.id], dt); });
    }
}

impl MongoDoc for PlanetSystem {
    fn get_id (&self) -> ObjectId {
        self.id
    }
}

impl PartialEq for PlanetSystem {
    fn eq(&self, other: &Self) -> bool {
        self.id == other.id
    }
}

impl Hash for PlanetSystem {
    fn hash<H: std::hash::Hasher>(&self, state: &mut H) {
        self.id.hash(state);
    }
}

impl Eq for PlanetSystem {}