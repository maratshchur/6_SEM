CREATE OR REPLACE TRIGGER trg_update_group_cval
AFTER INSERT OR UPDATE OR DELETE ON STUDENTS
FOR EACH ROW
BEGIN
    IF INSERTING THEN
        UPDATE GROUPS
        SET C_VAL = C_VAL + 1
        WHERE ID = :NEW.GROUP_ID;

    ELSIF UPDATING THEN
        IF :OLD.GROUP_ID != :NEW.GROUP_ID THEN
            UPDATE GROUPS
            SET C_VAL = C_VAL - 1
            WHERE ID = :OLD.GROUP_ID;

            UPDATE GROUPS
            SET C_VAL = C_VAL + 1
            WHERE ID = :NEW.GROUP_ID;
        END IF;

    ELSIF DELETING THEN
        UPDATE GROUPS
        SET C_VAL = C_VAL - 1
        WHERE ID = :OLD.GROUP_ID;
    END IF;
END;


-- -- Подготовка к тестированию
-- BEGIN
--     -- Очистка таблиц перед тестированием
--     DELETE FROM STUDENTS;
--     DELETE FROM GROUPS;
--     COMMIT;
-- END;


-- -- Создание группы для тестирования
-- BEGIN
--     INSERT INTO GROUPS (ID, NAME, C_VAL) VALUES (1, 'Group A', 0);
--     INSERT INTO GROUPS (ID, NAME, C_VAL) VALUES (2, 'Group B', 0);
--     COMMIT;
-- END;


-- -- Добавление студентов в группы
-- BEGIN
--     INSERT INTO STUDENTS (ID, NAME, GROUP_ID) VALUES (1, 'Alice', 1);
--     INSERT INTO STUDENTS (ID, NAME, GROUP_ID) VALUES (2, 'Bob', 1);
--     INSERT INTO STUDENTS (ID, NAME, GROUP_ID) VALUES (3, 'Charlie', 2);
--     COMMIT;
-- END;


-- -- Обновление группы студента
-- BEGIN
--     UPDATE STUDENTS SET GROUP_ID = 2 WHERE ID = 1;
--     COMMIT;
-- END;


-- -- Удаление студента
-- BEGIN
--     DELETE FROM STUDENTS WHERE ID = 2;
--     COMMIT;
-- END;


-- -- Проверка значений C_VAL в таблице GROUPS
-- BEGIN
--     FOR rec IN (SELECT * FROM GROUPS) LOOP
--         DBMS_OUTPUT.PUT_LINE('Group ID: ' || rec.ID || ', Name: ' || rec.NAME || ', C_VAL: ' || rec.C_VAL);
--     END LOOP;
-- END;
