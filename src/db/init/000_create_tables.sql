
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT uuid_generate_v4() UNIQUE,
    created timestamp DEFAULT now() NOT NULL,
    updated timestamp DEFAULT now() NOT NULL,
    name text NOT NULL,
    email text NOT NULL,
    is_admin boolean DEFAULT false NOT NULL,


    -- timestamp sanity check
    CHECK (created <= updated),

    PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS posts (
    id uuid DEFAULT uuid_generate_v4() UNIQUE,
    created timestamp DEFAULT now() NOT NULL,
    updated timestamp DEFAULT now() NOT NULL,

    created_by uuid NOT NULL,

    -- when does the post expire? Must be later than `created`
    expiration timestamp NOT NULL,
    
    -- traversing the comment chain. Both of these will be null
    -- if this is a top-level post. Otherwise, root_parent is the top,
    -- and direct_parent is one level up
    root_parent uuid,
    direct_parent uuid,

    contents text NOT NULL,
    approved boolean DEFAULT false,
    in_mod_queue boolean DEFAULT true,

    -- timestamp sanity checks
    CHECK (created <= updated),
    CHECK (created < expiration),

    

    PRIMARY KEY (id),
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);
