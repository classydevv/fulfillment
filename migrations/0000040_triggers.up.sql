CREATE TRIGGER update_updated_at_providers
    BEFORE UPDATE
    ON
        providers
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();