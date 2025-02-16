CREATE TABLE groups (
    id NUMBER NOT NULL,
    group_name VARCHAR2(20) NOT NULL,
    C_VAL NUMBER NOT NULL,
    CONSTRAINT group_id_pk PRIMARY KEY (id)
)
CREATE TABLE students (
    student_id NUMBER NOT NULL,
    student_name VARCHAR2(20) NOT NULL,
    gr_id NUMBER NOT NULL,
    CONSTRAINT student_id_pk PRIMARY KEY (student_id)
)

CREATE TABLE students_log (
    log_id NUMBER NOT NULL,
    action VARCHAR2(6) NOT NULL,
    new_student_id NUMBER,
    old_student_id NUMBER,
    new_student_name VARCHAR2(20),
    old_student_name VARCHAR2(20),
    new_gr_id NUMBER,
    old_gr_id NUMBER,
    action_date TIMESTAMP NOT NULL,
    CONSTRAINT log_id_pk PRIMARY KEY (log_id)
);

CREATE SEQUENCE STUDENT_ID_SEQ START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE GROUP_ID_SEQ START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE students_log_seq START WITH 1 INCREMENT BY 1;


CREATE OR REPLACE TRIGGER student_id_autoinc BEFORE
    INSERT ON students FOR EACH ROW 
BEGIN
    IF :NEW.STUDENT_ID IS NULL THEN
        :NEW.STUDENT_ID := STUDENT_ID_SEQ.NEXTVAL;
    END IF;
END;

CREATE OR REPLACE TRIGGER group_id_autoinc BEFORE
    INSERT ON groups FOR EACH ROW
BEGIN
    :NEW.id := GROUP_ID_SEQ.NEXTVAL;
END;

CREATE OR REPLACE TRIGGER cascade_delete_students
BEFORE DELETE ON groups
FOR EACH ROW
BEGIN
    DELETE FROM students
    WHERE gr_id = :old.id;
END;

CREATE OR REPLACE TRIGGER check_student_id_uniqueness
BEFORE INSERT ON students
FOR EACH ROW
DECLARE
    duplicate_count NUMBER;
BEGIN
    SELECT COUNT(*) INTO duplicate_count
    FROM students
    WHERE student_id = :new.student_id;

    IF duplicate_count > 0 THEN
        RAISE_APPLICATION_ERROR(-20001, 'Student ID must be unique ' || :new.student_id);
    END IF;
END;

CREATE OR REPLACE TRIGGER check_group_id_name_uniqueness
BEFORE INSERT OR UPDATE ON groups
FOR EACH ROW
DECLARE
    duplicate_id_count NUMBER;
    duplicate_name_count NUMBER;
BEGIN
    IF INSERTING OR UPDATING('ID') THEN
        SELECT COUNT(*) INTO duplicate_id_count
        FROM groups
        WHERE id = :new.id;

        IF duplicate_id_count > 0 THEN
            RAISE_APPLICATION_ERROR(-20001, 'group_id must be unique');
        END IF;
    END IF;

    IF INSERTING OR UPDATING('GROUP_NAME') THEN
        SELECT COUNT(*) INTO duplicate_name_count
        FROM groups
        WHERE LOWER(group_name) = LOWER(:new.group_name);

        IF duplicate_name_count > 0 THEN
            RAISE_APPLICATION_ERROR(-20001, 'group_name must be unique');
        END IF;
    END IF;
END;


CREATE OR REPLACE TRIGGER students_log_trigger
AFTER INSERT OR UPDATE OR DELETE ON students
FOR EACH ROW
BEGIN
    IF INSERTING THEN
        INSERT INTO students_log (log_id, action, new_student_id, new_student_name, new_gr_id, action_date)
        VALUES (students_log_seq.nextval, 'INSERT', :new.student_id, :new.student_name, :new.gr_id, SYSTIMESTAMP);
    
    ELSIF UPDATING THEN
        INSERT INTO students_log (log_id, action, old_student_id, new_student_id, old_student_name, new_student_name, old_gr_id, new_gr_id, action_date)
        VALUES (students_log_seq.nextval, 'UPDATE', :old.student_id, :new.student_id, :old.student_name, :new.student_name, :old.gr_id, :new.gr_id, SYSTIMESTAMP);
    
    ELSIF DELETING THEN
        INSERT INTO students_log (log_id, action, old_student_id, old_student_name, old_gr_id, action_date)
        VALUES (students_log_seq.nextval, 'DELETE', :old.student_id, :old.student_name, :old.gr_id, SYSTIMESTAMP);
    END IF;
END;


CREATE OR REPLACE PROCEDURE restore_students_state(
    p_timestamp TIMESTAMP,
    time_offset INTERVAL DAY TO SECOND DEFAULT NULL
) IS
    effective_time TIMESTAMP;
BEGIN
    IF time_offset IS NOT NULL THEN
        effective_time := SYSTIMESTAMP - time_offset;
    ELSE
        effective_time := p_timestamp;
    END IF;

    FOR lg IN (
        SELECT *
        FROM students_log
        WHERE action_date >= effective_time
        ORDER BY action_date DESC
    ) LOOP
        IF lg.action = 'INSERT' THEN
            DELETE FROM students
            WHERE student_id = lg.new_student_id;

        ELSIF lg.action = 'UPDATE' THEN
            UPDATE students
            SET student_name = lg.old_student_name,
                gr_id = lg.old_gr_id,
                student_id = lg.old_student_id
            WHERE student_id = lg.new_student_id;

        ELSIF lg.action = 'DELETE' THEN
            INSERT INTO students (student_id, student_name, gr_id)
            VALUES (lg.old_student_id, lg.old_student_name, lg.old_gr_id);
        END IF;
    END LOOP;
END;
/

CREATE OR REPLACE TRIGGER update_c_val_trigger AFTER
  INSERT OR UPDATE OR DELETE ON students FOR EACH ROW
DECLARE
  new_group_id NUMBER;
  old_group_id NUMBER;
BEGIN
  IF inserting THEN
    UPDATE groups
    SET
      c_val = c_val + 1
    WHERE
      id = :new.gr_id;
  ELSIF updating THEN
    new_group_id := :new.gr_id;
    old_group_id := :old.gr_id;
    IF new_group_id <> old_group_id THEN
      UPDATE groups
      SET
        c_val = c_val + 1
      WHERE
        id = new_group_id;
      UPDATE groups
      SET
        c_val = c_val - 1
      WHERE
        id = old_group_id;
    END IF;
  ELSIF deleting THEN
    UPDATE groups
    SET
      c_val = c_val - 1
    WHERE
      id = :old.gr_id;
  END IF;
EXCEPTION
  WHEN OTHERS THEN
    dbms_output.put_line('Error updating c_val ' || sqlerrm);
END;