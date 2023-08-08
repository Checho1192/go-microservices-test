CREATE TABLE MEASUREMENTS(
    id UUID PRIMARY KEY,
    meter_id INT,
    active_energy DECIMAL,
    reactive_energy DECIMAL,
    capacitive_reactive DECIMAL,
    solar DECIMAL,
    date TIMESTAMP
);
--- For importing data from CSV files, use the following command:
--- COPY measurements FROM '/path/to/csv/file' DELIMITER ',' CSV HEADER;
--- You can also import data using Dbeaver import data tool:
--- Right click on table -> Import Data -> Select CSV file -> Next -> Next -> Next -> Next -> Finish
```