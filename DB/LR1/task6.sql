CREATE OR REPLACE FUNCTION calculate_total_reward(p_monthly_salary NUMBER, p_bonus_percentage NUMBER)
RETURN NUMBER IS
    v_total_reward NUMBER;
BEGIN
    IF p_monthly_salary < 0 OR p_bonus_percentage < 0 THEN
        RAISE_APPLICATION_ERROR(-20001, 'Invalid input values.');
    END IF;

    v_total_reward := (1 + p_bonus_percentage / 100) * 12 * p_monthly_salary;
    RETURN v_total_reward;
END;

DECLARE
    v_result NUMBER;
BEGIN
    v_result := calculate_total_reward(10, 5);
    
    DBMS_OUTPUT.PUT_LINE('Общая награда: ' || v_result);
EXCEPTION
    WHEN OTHERS THEN
        DBMS_OUTPUT.PUT_LINE('Ошибка: ' || SQLERRM);
END;