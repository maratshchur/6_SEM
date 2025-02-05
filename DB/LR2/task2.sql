CREATE OR REPLACE TRIGGER trg_groups_before_insert
BEFORE INSERT ON GROUPS
FOR EACH ROW
DECLARE
    v_count NUMBER;
BEGIN
    SELECT COUNT(*)
    INTO v_count
    FROM GROUPS
    WHERE ID = :NEW.ID;

    IF v_count > 0 THEN
        RAISE_APPLICATION_ERROR(-20001, 'ID must be unique.');
    END IF;

    IF :NEW.ID IS NULL THEN
        SELECT NVL(MAX(ID), 0) + 1
        INTO :NEW.ID
        FROM GROUPS;
    END IF;

    SELECT COUNT(*)
    INTO v_count
    FROM GROUPS
    WHERE NAME = :NEW.NAME;

    IF v_count > 0 THEN
        RAISE_APPLICATION_ERROR(-20002, 'Group name must be unique.');
    END IF;
END;


CREATE OR REPLACE TRIGGER trg_students_before_insert
BEFORE INSERT ON STUDENTS
FOR EACH ROW
DECLARE
    v_count NUMBER;
BEGIN
    SELECT COUNT(*)
    INTO v_count
    FROM STUDENTS
    WHERE ID = :NEW.ID;

    IF v_count > 0 THEN
        RAISE_APPLICATION_ERROR(-20001, 'ID must be unique.');
    END IF;

    IF :NEW.ID IS NULL THEN
        SELECT NVL(MAX(ID), 0) + 1
        INTO :NEW.ID
        FROM STUDENTS;
    END IF;
END;

-- BEGIN
--     INSERT INTO GROUPS (NAME, C_VAL) VALUES ('Group A', 1);
    
--     INSERT INTO GROUPS (ID, NAME, C_VAL) VALUES (2, 'Group B', 0);
    
--     BEGIN
--         INSERT INTO GROUPS (ID, NAME, C_VAL) VALUES (2, 'Group C', 0);
--     EXCEPTION
--         WHEN OTHERS THEN
--             DBMS_OUTPUT.PUT_LINE(SQLERRM);
--     END;

--     BEGIN
--         INSERT INTO GROUPS (ID, NAME, C_VAL) VALUES (3, 'Group B', 0);
--     EXCEPTION
--         WHEN OTHERS THEN
--             DBMS_OUTPUT.PUT_LINE(SQLERRM);
--     END;

--     INSERT INTO GROUPS (NAME, C_VAL) VALUES ('Group D', 2);
    
--     COMMIT;
-- END;


-- BEGIN
--     INSERT INTO STUDENTS (ID, NAME, GROUP_ID) VALUES (1, 'Alice', 1);
    
--     INSERT INTO STUDENTS (NAME, GROUP_ID) VALUES ('Bob', 1);
    
--     BEGIN
--         INSERT INTO STUDENTS (ID, NAME, GROUP_ID) VALUES (1, 'Charlie', 1);
--     EXCEPTION
--         WHEN OTHERS THEN
--             DBMS_OUTPUT.PUT_LINE('Error: ' || SQLERRM);
--     END;

--     INSERT INTO STUDENTS (ID, NAME, GROUP_ID) VALUES (3, 'David', 1);
    
--     FOR rec IN (SELECT * FROM STUDENTS) LOOP
--         DBMS_OUTPUT.PUT_LINE('ID: ' || rec.ID || ', Name: ' || rec.NAME || ', Group ID: ' || rec.GROUP_ID);
--     END LOOP;
    
--          COMMIT;


-- END;
