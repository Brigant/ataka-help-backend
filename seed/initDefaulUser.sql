-- change salt to your own value
INSERT INTO public.users
(email, password, firs_name, last_name)
VALUES('admin@example.com', encode(sha256(concat('super-password', 'salt')::bytea), 'hex'), 'admin', 'adminovna');