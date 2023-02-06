

-- Set timestamps on the rows. These triggers have the effect of setting `created` on insert, and `updated` any time the
-- row is modified

CREATE FUNCTION touch_created() RETURNS trigger AS $$
    BEGIN
        NEW.created = now();
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION touch_updated() RETURNS trigger AS $$
    BEGIN
        NEW.updated = now();
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER user_set_created BEFORE INSERT ON users
    FOR EACH ROW EXECUTE FUNCTION touch_created();

CREATE TRIGGER user_set_updated BEFORE INSERT OR UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION touch_updated();

CREATE TRIGGER post_set_created BEFORE INSERT ON posts
    FOR EACH ROW EXECUTE FUNCTION touch_created();

CREATE TRIGGER post_set_updated BEFORE INSERT OR UPDATE ON posts
    FOR EACH ROW EXECUTE FUNCTION touch_updated();