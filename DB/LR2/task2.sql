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
--     -- Вставка группы 1
--     INSERT INTO GROUPS (NAME, C_VAL) VALUES ('Group A', 1);
    
--     -- Вставка группы 2
--     INSERT INTO GROUPS (ID, NAME, C_VAL) VALUES (2, 'Group B', 0);
    
--     -- Попытка вставить группу с существующим ID
--     BEGIN
--         INSERT INTO GROUPS (ID, NAME, C_VAL) VALUES (2, 'Group C', 0);
--     EXCEPTION
--         WHEN OTHERS THEN
--             DBMS_OUTPUT.PUT_LINE(SQLERRM);
--     END;

--     -- Попытка вставить группу с существующим NAME
--     BEGIN
--         INSERT INTO GROUPS (ID, NAME, C_VAL) VALUES (3, 'Group B', 0);
--     EXCEPTION
--         WHEN OTHERS THEN
--             DBMS_OUTPUT.PUT_LINE(SQLERRM);
--     END;

--     -- Вставка группы без ID (должен сгенерироваться автоинкрементный ID)
--     INSERT INTO GROUPS (NAME, C_VAL) VALUES ('Group D', 2);
    
--     -- Проверка результатов
--     COMMIT;
-- END;


-- -- Тестирование триггера для таблицы STUDENTS
-- BEGIN
--     -- Вставка студента с уникальным ID
--     INSERT INTO STUDENTS (ID, NAME, GROUP_ID) VALUES (1, 'Alice', 1);
    
--     -- Вставка студента без указания ID (должен сгенерироваться автоинкрементный ID)
--     INSERT INTO STUDENTS (NAME, GROUP_ID) VALUES ('Bob', 1);
    
--     -- Попытка вставить студента с существующим ID
--     BEGIN
--         INSERT INTO STUDENTS (ID, NAME, GROUP_ID) VALUES (1, 'Charlie', 1);
--     EXCEPTION
--         WHEN OTHERS THEN
--             DBMS_OUTPUT.PUT_LINE('Error: ' || SQLERRM);
--     END;

--     -- Вставка студента с уникальным ID
--     INSERT INTO STUDENTS (ID, NAME, GROUP_ID) VALUES (3, 'David', 1);
    
--     -- Проверка текущих записей в таблице STUDENTS
--     FOR rec IN (SELECT * FROM STUDENTS) LOOP
--         DBMS_OUTPUT.PUT_LINE('ID: ' || rec.ID || ', Name: ' || rec.NAME || ', Group ID: ' || rec.GROUP_ID);
--     END LOOP;
    
--          COMMIT;


-- END;
