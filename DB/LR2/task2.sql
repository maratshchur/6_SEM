CREATE OR REPLACE TRIGGER trg_unique_students_id
BEFORE INSERT ON STUDENTS
FOR EACH ROW
DECLARE
    v_count NUMBER;
BEGIN
    SELECT COUNT(*) INTO v_count FROM STUDENTS WHERE ID = :NEW.ID;
    IF v_count > 0 THEN
        RAISE_APPLICATION_ERROR(-20001, 'ID студента должен быть уникальным');
    END IF;
END;