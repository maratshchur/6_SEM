CREATE OR REPLACE TRIGGER trg_students_before_delete
BEFORE DELETE ON GROUPS
FOR EACH ROW
BEGIN
    DELETE FROM STUDENTS WHERE GROUP_ID = :OLD.ID;
END;

-- BEGIN
--     DELETE FROM STUDENTS;
--     DELETE FROM GROUPS;
--     COMMIT;
-- END;
-- /

-- BEGIN
--     INSERT INTO GROUPS (ID, NAME, C_VAL) VALUES (1, 'Group A', 0);
--     INSERT INTO STUDENTS (ID, NAME, GROUP_ID) VALUES (1, 'Alice', 1);
--     INSERT INTO STUDENTS (ID, NAME, GROUP_ID) VALUES (2, 'Bob', 1);
    
--     COMMIT;
-- END;
-- /

-- BEGIN
--     DELETE FROM GROUPS WHERE ID = 1;
    
--     FOR rec IN (SELECT * FROM STUDENTS) LOOP
--         DBMS_OUTPUT.PUT_LINE('ID: ' || rec.ID || ', Name: ' || rec.NAME || ', Group ID: ' || rec.GROUP_ID);
--     END LOOP;

--     COMMIT;
-- END;
-- /