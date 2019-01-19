--
-- PostgreSQL database dump
--

-- Dumped from database version 10.2
-- Dumped by pg_dump version 10.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE accounts (
    id bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    deleted_at timestamp without time zone,
    status smallint NOT NULL,
    name character varying NOT NULL,
    phone character varying,
    country_id integer
);


--
-- Name: accounts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE accounts_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: accounts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE accounts_id_seq OWNED BY accounts.id;


--
-- Name: amenities; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE amenities (
    id bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    type character varying NOT NULL,
    name character varying,
    property_id integer,
    unit_id integer,
    CONSTRAINT amenities_check CHECK ((((type)::text <> 'custom'::text) OR (name IS NOT NULL))),
    CONSTRAINT amenities_check1 CHECK (((property_id IS NOT NULL) OR (unit_id IS NOT NULL)))
);


--
-- Name: amenities_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE amenities_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: amenities_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE amenities_id_seq OWNED BY amenities.id;


--
-- Name: bookings; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE bookings (
    id bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    user_id integer NOT NULL,
    unit_id integer NOT NULL,
    check_in date NOT NULL,
    check_out date NOT NULL,
    message text,
    night_price_cents integer NOT NULL,
    nights integer NOT NULL,
    service_fee_cents integer NOT NULL,
    total_cents integer NOT NULL,
    CONSTRAINT bookings_nights_check CHECK ((nights > 0))
);


--
-- Name: bookings_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE bookings_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: bookings_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE bookings_id_seq OWNED BY bookings.id;


--
-- Name: cities; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE cities (
    id bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    deleted_at timestamp without time zone,
    name character varying NOT NULL,
    slug character varying NOT NULL,
    country_id integer NOT NULL
);


--
-- Name: cities_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE cities_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: cities_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE cities_id_seq OWNED BY cities.id;


--
-- Name: countries; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE countries (
    id bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    deleted_at timestamp without time zone,
    name character varying NOT NULL
);


--
-- Name: countries_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE countries_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: countries_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE countries_id_seq OWNED BY countries.id;


--
-- Name: images; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE images (
    id bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    name character varying NOT NULL,
    type character varying NOT NULL,
    size integer NOT NULL,
    url character varying NOT NULL,
    property_id integer,
    unit_id integer,
    uploaded boolean DEFAULT false,
    "position" smallint,
    CONSTRAINT images_check CHECK (((property_id IS NOT NULL) OR (unit_id IS NOT NULL)))
);


--
-- Name: images_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE images_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: images_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE images_id_seq OWNED BY images.id;


--
-- Name: properties; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE properties (
    id bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    deleted_at timestamp without time zone,
    account_id integer NOT NULL,
    published_at timestamp without time zone,
    ever_published boolean DEFAULT false NOT NULL,
    name character varying,
    slug character varying,
    address1 character varying,
    address2 character varying,
    address3 character varying,
    city_id integer,
    postal_code character varying,
    country_id integer,
    latitude numeric,
    longitude numeric,
    minimum_stay smallint,
    deposit character varying,
    cleaning character varying,
    nearest_airport character varying,
    nearest_subway character varying,
    nearby_locations character varying,
    overview character varying
);


--
-- Name: properties_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE properties_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: properties_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE properties_id_seq OWNED BY properties.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


--
-- Name: units; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE units (
    id bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    deleted_at timestamp without time zone,
    property_id integer NOT NULL,
    published_at timestamp without time zone,
    ever_published boolean DEFAULT false NOT NULL,
    name character varying NOT NULL,
    slug character varying NOT NULL,
    bedrooms integer NOT NULL,
    bathrooms integer NOT NULL,
    size_m2 integer NOT NULL,
    max_occupancy integer,
    count integer NOT NULL,
    base_price_cents integer,
    one_night_price_cents integer,
    one_week_price_cents integer,
    three_months_price_cents integer,
    six_months_price_cents integer,
    twelve_months_price_cents integer,
    overview text,
    floor_plan_image_id integer
);


--
-- Name: units_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE units_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: units_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE units_id_seq OWNED BY units.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE users (
    id bigint NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    deleted_at timestamp without time zone,
    account_id integer,
    email character varying NOT NULL,
    name character varying,
    facebook_id character varying,
    encrypted_password character varying,
    last_signed_in timestamp without time zone
);


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE users_id_seq OWNED BY users.id;


--
-- Name: accounts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY accounts ALTER COLUMN id SET DEFAULT nextval('accounts_id_seq'::regclass);


--
-- Name: amenities id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY amenities ALTER COLUMN id SET DEFAULT nextval('amenities_id_seq'::regclass);


--
-- Name: bookings id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY bookings ALTER COLUMN id SET DEFAULT nextval('bookings_id_seq'::regclass);


--
-- Name: cities id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY cities ALTER COLUMN id SET DEFAULT nextval('cities_id_seq'::regclass);


--
-- Name: countries id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY countries ALTER COLUMN id SET DEFAULT nextval('countries_id_seq'::regclass);


--
-- Name: images id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY images ALTER COLUMN id SET DEFAULT nextval('images_id_seq'::regclass);


--
-- Name: properties id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY properties ALTER COLUMN id SET DEFAULT nextval('properties_id_seq'::regclass);


--
-- Name: units id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY units ALTER COLUMN id SET DEFAULT nextval('units_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT nextval('users_id_seq'::regclass);


--
-- Name: accounts accounts_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY accounts
    ADD CONSTRAINT accounts_name_key UNIQUE (name);


--
-- Name: accounts accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);


--
-- Name: amenities amenities_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY amenities
    ADD CONSTRAINT amenities_pkey PRIMARY KEY (id);


--
-- Name: bookings bookings_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY bookings
    ADD CONSTRAINT bookings_pkey PRIMARY KEY (id);


--
-- Name: cities cities_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY cities
    ADD CONSTRAINT cities_name_key UNIQUE (name);


--
-- Name: cities cities_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY cities
    ADD CONSTRAINT cities_pkey PRIMARY KEY (id);


--
-- Name: cities cities_slug_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY cities
    ADD CONSTRAINT cities_slug_key UNIQUE (slug);


--
-- Name: countries countries_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY countries
    ADD CONSTRAINT countries_name_key UNIQUE (name);


--
-- Name: countries countries_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY countries
    ADD CONSTRAINT countries_pkey PRIMARY KEY (id);


--
-- Name: images images_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY images
    ADD CONSTRAINT images_pkey PRIMARY KEY (id);


--
-- Name: properties properties_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY properties
    ADD CONSTRAINT properties_name_key UNIQUE (name);


--
-- Name: properties properties_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY properties
    ADD CONSTRAINT properties_pkey PRIMARY KEY (id);


--
-- Name: properties properties_slug_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY properties
    ADD CONSTRAINT properties_slug_key UNIQUE (slug);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: units units_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY units
    ADD CONSTRAINT units_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: accounts_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX accounts_deleted_at ON accounts USING btree (deleted_at);


--
-- Name: amenities_name_uniq_by_property; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX amenities_name_uniq_by_property ON amenities USING btree (property_id, name);


--
-- Name: amenities_name_uniq_by_unit; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX amenities_name_uniq_by_unit ON amenities USING btree (unit_id, name);


--
-- Name: bookings_unit_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX bookings_unit_id ON bookings USING btree (unit_id);


--
-- Name: bookings_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX bookings_user_id ON bookings USING btree (user_id);


--
-- Name: images_property_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX images_property_id ON images USING btree (property_id);


--
-- Name: images_unit_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX images_unit_id ON images USING btree (unit_id);


--
-- Name: properties_account_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX properties_account_id ON properties USING btree (account_id);


--
-- Name: properties_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX properties_deleted_at ON properties USING btree (deleted_at);


--
-- Name: properties_published_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX properties_published_at ON properties USING btree (published_at);


--
-- Name: units_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX units_deleted_at ON units USING btree (deleted_at);


--
-- Name: units_property_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX units_property_id ON units USING btree (property_id);


--
-- Name: units_property_id_slug; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX units_property_id_slug ON units USING btree (property_id, slug);


--
-- Name: units_published_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX units_published_at ON units USING btree (published_at);


--
-- Name: users_deleted_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX users_deleted_at ON users USING btree (deleted_at);


--
-- Name: users_email; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX users_email ON users USING btree (email);


--
-- Name: accounts accounts_country_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY accounts
    ADD CONSTRAINT accounts_country_id_fkey FOREIGN KEY (country_id) REFERENCES countries(id);


--
-- Name: amenities amenities_property_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY amenities
    ADD CONSTRAINT amenities_property_id_fkey FOREIGN KEY (property_id) REFERENCES properties(id);


--
-- Name: amenities amenities_unit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY amenities
    ADD CONSTRAINT amenities_unit_id_fkey FOREIGN KEY (unit_id) REFERENCES units(id);


--
-- Name: bookings bookings_unit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY bookings
    ADD CONSTRAINT bookings_unit_id_fkey FOREIGN KEY (unit_id) REFERENCES units(id);


--
-- Name: bookings bookings_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY bookings
    ADD CONSTRAINT bookings_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id);


--
-- Name: cities cities_country_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY cities
    ADD CONSTRAINT cities_country_id_fkey FOREIGN KEY (country_id) REFERENCES countries(id);


--
-- Name: images images_property_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY images
    ADD CONSTRAINT images_property_id_fkey FOREIGN KEY (property_id) REFERENCES properties(id);


--
-- Name: images images_unit_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY images
    ADD CONSTRAINT images_unit_id_fkey FOREIGN KEY (unit_id) REFERENCES units(id);


--
-- Name: properties properties_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY properties
    ADD CONSTRAINT properties_account_id_fkey FOREIGN KEY (account_id) REFERENCES accounts(id);


--
-- Name: properties properties_city_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY properties
    ADD CONSTRAINT properties_city_id_fkey FOREIGN KEY (city_id) REFERENCES cities(id);


--
-- Name: properties properties_country_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY properties
    ADD CONSTRAINT properties_country_id_fkey FOREIGN KEY (country_id) REFERENCES countries(id);


--
-- Name: units units_floor_plan_image_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY units
    ADD CONSTRAINT units_floor_plan_image_id_fkey FOREIGN KEY (floor_plan_image_id) REFERENCES images(id);


--
-- Name: units units_property_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY units
    ADD CONSTRAINT units_property_id_fkey FOREIGN KEY (property_id) REFERENCES properties(id);


--
-- Name: users users_account_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_account_id_fkey FOREIGN KEY (account_id) REFERENCES accounts(id);


--
-- PostgreSQL database dump complete
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 10.2
-- Dumped by pg_dump version 10.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

SET search_path = public, pg_catalog;

--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY schema_migrations (version, dirty) FROM stdin;
9	f
\.


--
-- PostgreSQL database dump complete
--

