package org.proj.db.codecs;

import org.bson.BsonReader;
import org.bson.BsonWriter;
import org.bson.codecs.Codec;
import org.bson.codecs.DecoderContext;
import org.bson.codecs.EncoderContext;
import org.proj.db.codecs.primitive.ArrayCodec;
import org.proj.game.body.Planet;
import org.proj.game.PlanetarySystem;
import org.proj.game.body.Sun;

public class PlanetarySystemCodec implements Codec<PlanetarySystem> {
    final public static PlanetarySystemCodec INSTANCE = new PlanetarySystemCodec();

    final private ArrayCodec<Sun> SUN = new ArrayCodec<>(SunCodec.INSTANCE);
    final private ArrayCodec<Planet> PLANET = new ArrayCodec<>(PlanetCodec.INSTANCE);

    private PlanetarySystemCodec () {}

    @Override
    public PlanetarySystem decode (BsonReader reader, DecoderContext context) {
        reader.readStartDocument();
        reader.readBsonType(); reader.skipName(); reader.skipValue();

        reader.readName("suns");
        Sun[] suns = SUN.decode(reader, context);

        reader.readName("planets");
        Planet[] planets = PLANET.decode(reader, context);

        reader.readEndDocument();
        return new PlanetarySystem(suns, planets);
    }

    @Override
    public void encode (BsonWriter writer, PlanetarySystem value, EncoderContext context) {
        writer.writeStartDocument();

        writer.writeName("suns");
        SUN.encode(writer, value.getSuns(), context);

        writer.writeName("planets");
        PLANET.encode(writer, value.getPlanets(), context);

        writer.writeEndDocument();
    }

    @Override
    public Class<PlanetarySystem> getEncoderClass() {
        return PlanetarySystem.class;
    }
}