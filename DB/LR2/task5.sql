CREATE OR REPLACE PROCEDURE restore_students(
    p_time TIMESTAMP,
    p_offset INTERVAL DAY TO SECOND DEFAULT NULL
) AS
    v_adjusted_time TIMESTAMP;
    v_students_count NUMBER;
BEGIN
    -- Если указано временное смещение, вычисляем новое время
    IF p_offset IS NOT NULL THEN
        v_adjusted_time := p_time + p_offset;
    ELSE
        v_adjusted_time := p_time;
    END IF;

    -- Очищаем текущую таблицу STUDENTS
    DELETE FROM STUDENTS;

    -- Восстанавливаем данные из журнала на указанное время
    FOR rec IN (
        SELECT STUDENT_ID, ACTION, NEW_NAME, NEW_GROUP_ID, OLD_NAME, OLD_GROUP_ID
        FROM STUDENTS_LOG
        WHERE ACTION_TIMESTAMP <= v_adjusted_time
        ORDER BY ACTION_TIMESTAMP DESC
    ) LOOP
        IF rec.ACTION = 'INSERT' THEN
            INSERT INTO STUDENTS (ID, NAME, GROUP_ID) 
            VALUES (rec.STUDENT_ID, rec.NEW_NAME, rec.NEW_GROUP_ID);
        ELSIF rec.ACTION = 'UPDATE' THEN
            INSERT INTO STUDENTS (ID, NAME, GROUP_ID) 
            VALUES (rec.STUDENT_ID, rec.OLD_NAME, rec.OLD_GROUP_ID);
        ELSIF rec.ACTION = 'DELETE' THEN
            -- Если запись была удалена, её следует восстановить
            INSERT INTO STUDENTS (ID, NAME, GROUP_ID) 
            VALUES (rec.STUDENT_ID, rec.OLD_NAME, rec.OLD_GROUP_ID);
        END IF;
    END LOOP;

    COMMIT;

    -- Проверка восстановленных данных
    SELECT COUNT(*) INTO v_students_count FROM STUDENTS;
    DBMS_OUTPUT.PUT_LINE('Restored ' || v_students_count || ' students to the state as of ' || v_adjusted_time);
END;

-- -- Пример использования процедуры
-- DECLARE
--     restore_time TIMESTAMP := SYSTIMESTAMP - INTERVAL '1' HOUR; -- Восстановление данных на 1 час назад
-- BEGIN
--     restore_students(restore_time);
-- END;
-- /