use actix_web::{Responder, web, HttpRequest, post, get, HttpResponse};
use jsonwebtoken::{decode, DecodingKey, Validation};
use mongodb::{bson::{doc}};
use serde_json::{json, Value};
use strum::IntoEnumIterator;
use crate::{DATABASE, Resource, PLAYERS, Player, PlayerToken, PlayerTokenLoged, CURRENT_LOGGER, Logger, get_auth_token, decode_token};

// OUT API
#[get("/status")]
pub async fn status () -> impl Responder {
    let ping = match DATABASE.get().unwrap().run_command(doc! { "ping": 1u32 }, None).await {
        Err(x) => { tokio::spawn(CURRENT_LOGGER.log_error(x)); false },
        Ok(x) => match x.get_f64("ok") {
            Ok(x) => {
                let status = x == 1.;
                tokio::spawn(CURRENT_LOGGER.log_info(format!("Status is {x}")));
                status
            },
            Err(x) => { tokio::spawn(CURRENT_LOGGER.log_error(x)); false },
        }
    };

    let json = json!({
        "running": true,
        "database": ping
    });

    web::Json(json)
}

#[get("/resources")]
pub async fn resources () -> impl Responder {
    let resources = Resource::iter()
        .map(|x: Resource| json!({
            "name": format!("{:?}", x),
            "size": x.get_size(),
            "type": x.get_type()
        })).collect::<Vec<_>>();

    web::Json(json!({
        "resources": resources
    }))
}

// IN API
#[post("/player/signup")]
pub async fn new_user (_: HttpRequest, body: web::Json<u64>) -> HttpResponse {
    // TODO INTERNAL IP ONLY
    match PLAYERS.insert_one(Player::new(PlayerToken::Unloged(body.0), format!("todo"))).await {
        Ok(_) => HttpResponse::Ok().finish(),
        Err(x) => { 
            tokio::spawn(CURRENT_LOGGER.log_error(format!("{x}")));
            HttpResponse::InternalServerError().body(format!("{x}"))
        }
    }
}

#[post("/player/signin")]
pub async fn user_login (req: HttpRequest) -> HttpResponse {
    match decode_token(&req) {
        Ok((body, token)) => {
            let query = bson::to_document(&PlayerToken::Unloged(token.claims.id)).unwrap();
            let update = bson::to_document(&PlayerToken::Loged(body)).unwrap();
        
            match PLAYERS.update_one(doc! { "token": query }, doc! { "$set": { "token": update } }).await {
                Ok(Some(_)) => HttpResponse::Ok().finish(),
                Ok(None) => HttpResponse::BadRequest().body("No matching player found"),
                Err(e) => { 
                    tokio::spawn(CURRENT_LOGGER.log_error(format!("{e}"))); 
                    return HttpResponse::InternalServerError().body(format!("{e}"))
                }
            }
        },

        Err(e) => HttpResponse::BadRequest().body(format!("{e}"))
    }
}

#[post("/player/signout")]
pub async fn user_logout (req: HttpRequest) -> HttpResponse {
    if let Ok((string, token)) = decode_token(&req) {
        let update = bson::to_document(&PlayerToken::Unloged(token.claims.id)).unwrap();
        let query = bson::to_document(&PlayerToken::Loged(string)).unwrap();
    
        return match PLAYERS.update_one(doc! { "token": query }, doc! { "$set": { "token": update } }).await {
            Ok(Some(_)) => HttpResponse::Ok().finish(),
            Ok(None) => HttpResponse::BadRequest().body("No matching player found"),
            Err(e) => { 
                tokio::spawn(CURRENT_LOGGER.log_error(format!("{e}")));
                HttpResponse::InternalServerError().body(format!("{e}"))
            }
        }
    }

    HttpResponse::BadRequest().body("No authorization key found")
}