CREATE TABLE ais_data (
    -- Flattened fields from vesselParticulars
    vesselName TEXT,
    callSign TEXT,
    imoNumber TEXT,
    flag TEXT,
    vesselLength REAL,
    vesselBreadth REAL,
    vesselDepth REAL,
    vesselType TEXT,
    grossTonnage REAL,
    netTonnage REAL,
    deadweight REAL,
    mmsiNumber TEXT,
    yearBuilt TEXT,

    -- Root level fields
    latitude REAL,
    longitude REAL,
    latitudeDegrees REAL,
    longitudeDegrees REAL,
    speed REAL,
    course REAL,
    heading REAL,
    dimA INTEGER,
    dimB INTEGER,
    timeStamp TEXT
);

