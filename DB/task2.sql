DECLARE
    v_id NUMBER;
    v_val NUMBER;
BEGIN
    FOR i IN 1..10000 LOOP
        v_id := i;
        v_val := TRUNC(DBMS_RANDOM.VALUE(1, 1000));
        INSERT INTO MyTable (id, val) VALUES (v_id, v_val);
    END LOOP;
    COMMIT;
END;