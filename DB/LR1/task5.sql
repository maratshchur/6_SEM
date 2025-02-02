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

call insert_into_mytable(11000, 17);

select val from MyTable
where id = 11000;

call update_mytable(11000, 18);


select val from MyTable
where id = 11000;

call delete_from_mytable(11000);


select val from MyTable
where id = 11000;
