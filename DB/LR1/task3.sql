CREATE OR REPLACE FUNCTION check_even_odd RETURN VARCHAR2 IS
    v_even_count NUMBER;
    v_odd_count NUMBER;
BEGIN
    SELECT COUNT(*) INTO v_even_count FROM MyTable WHERE MOD(val, 2) = 0;
    SELECT COUNT(*) INTO v_odd_count FROM MyTable WHERE MOD(val, 2) = 1;

    IF v_even_count > v_odd_count THEN
        RETURN 'TRUE';
    ELSIF v_odd_count > v_even_count THEN
        RETURN 'FALSE';
    ELSE
        RETURN 'EQUAL';
    END IF;
END;

SELECT check_even_odd() AS result FROM dual;