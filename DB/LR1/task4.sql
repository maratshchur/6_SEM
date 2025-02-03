CREATE OR REPLACE FUNCTION generate_insert_command(p_id NUMBER) RETURN VARCHAR2 IS
    v_val NUMBER;
    v_insert_command VARCHAR2(1000);
BEGIN
    SELECT val INTO v_val FROM MyTable WHERE id = p_id;
    v_insert_command := 'INSERT INTO MyTable (id, val) VALUES (' || p_id || ', ' || v_val || ');';
    RETURN v_insert_command;
EXCEPTION
    WHEN NO_DATA_FOUND THEN
        RETURN 'ID not found.';
END;

DECLARE
    insert_command VARCHAR2(100); -- или другой подходящий размер
BEGIN
    insert_command := generate_insert_command(100000);
    
    DBMS_OUTPUT.PUT_LINE(insert_command);
END;
