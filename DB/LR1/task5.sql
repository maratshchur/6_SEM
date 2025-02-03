CREATE OR REPLACE PROCEDURE insert_into_mytable(p_id NUMBER, p_val NUMBER) IS
BEGIN
    INSERT INTO MyTable (id, val) VALUES (p_id, p_val);
    COMMIT;
END;

CREATE OR REPLACE PROCEDURE update_mytable(p_id NUMBER, p_val NUMBER) IS
BEGIN
    UPDATE MyTable SET val = p_val WHERE id = p_id;
    COMMIT;
END;

CREATE OR REPLACE PROCEDURE delete_from_mytable(p_id NUMBER) IS
BEGIN
    DELETE FROM MyTable WHERE id = p_id;
    COMMIT;
END;

DECLARE
    v_val NUMBER; 
BEGIN
    insert_into_mytable(11000, 17);

    BEGIN
        SELECT val INTO v_val FROM MyTable WHERE id = 11000;
        DBMS_OUTPUT.PUT_LINE('Значение после вставки: ' || v_val);
    EXCEPTION
        WHEN NO_DATA_FOUND THEN
            DBMS_OUTPUT.PUT_LINE('Запись не найдена после вставки.');
    END;

    update_mytable(11000, 18);

    BEGIN
        SELECT val INTO v_val FROM MyTable WHERE id = 11000;
        DBMS_OUTPUT.PUT_LINE('Значение после обновления: ' || v_val);
    EXCEPTION
        WHEN NO_DATA_FOUND THEN
            DBMS_OUTPUT.PUT_LINE('Запись не найдена после обновления.');
    END;

    delete_from_mytable(11000);

    BEGIN
        SELECT val INTO v_val FROM MyTable WHERE id = 11000;
        DBMS_OUTPUT.PUT_LINE('Запись все еще существует после удаления: ' || v_val);
    EXCEPTION
        WHEN NO_DATA_FOUND THEN
            DBMS_OUTPUT.PUT_LINE('Запись успешно удалена.');
    END;

END;
