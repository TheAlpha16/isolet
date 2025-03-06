-- Create categories
INSERT INTO categories (category_id, category_name) VALUES
    (1, 'Miscellaneous'),
    (2, 'Reversing'),
    (3, 'Pwn'),
    (4, 'Linux'),
    (5, 'Web'),
    (6, 'Crypto'),
    (7, 'Forensics'),
    (8, 'OSINT'),
    (9, 'Steganography');

-- Add challenge data
INSERT INTO challenges (
    chall_id, chall_name, category_id, 
    prompt, flag, type, 
    points, files, requirements, 
    author, tags, links,
    deployment, port, subd
) VALUES 
(1, 'challenge_1', 2, 'Solve challenge_1!', 'isolet-dev{random_flag_1}', 'on-demand', 200, '{}', '{}', 'CyberKid', '{python,misc}', '{}', 'http', 443, 'subd_1'),
(2, 'challenge_2', 1, 'Solve challenge_2!', 'isolet-dev{random_flag_2}', 'static', 300, '{file_2.bin}', '{}', 'TheAlpha', '{linux,df,web}', '{}', 'http', 80, ''),
(3, 'challenge_3', 2, 'Solve challenge_3!', 'isolet-dev{random_flag_3}', 'dynamic', 200, '{file_3.bin}', '{}', 'Naughtyb0y', '{df,misc}', '{}', 'http', 80, ''),
(4, 'challenge_4', 8, 'Solve challenge_4!', 'isolet-dev{random_flag_4}', 'dynamic', 500, '{file_4.bin}', '{}', 'CyberKid', '{heap,df}', '{}', 'nc', 30086, ''),
(5, 'challenge_5', 2, 'Solve challenge_5!', 'isolet-dev{random_flag_5}', 'static', 300, '{file_5.bin}', '{}', 'Hackerman', '{heap}', '{}', 'http', 80, ''),
(6, 'challenge_6', 9, 'Solve challenge_6!', 'isolet-dev{random_flag_6}', 'dynamic', 500, '{file_6.bin}', '{}', 'e4stw1nd', '{misc,linux,web}', '{}', 'nc', 31000, ''),
(7, 'challenge_7', 2, 'Solve challenge_7!', 'isolet-dev{random_flag_7}', 'static', 500, '{file_7.bin}', '{}', 'CyberKid', '{forensics,reversing}', '{}', 'http', 80, ''),
(8, 'challenge_8', 9, 'Solve challenge_8!', 'isolet-dev{random_flag_8}', 'static', 400, '{file_8.bin}', '{}', 'Hackerman', '{web,heap}', '{}', 'http', 80, ''),
(9, 'challenge_9', 6, 'Solve challenge_9!', 'isolet-dev{random_flag_9}', 'static', 300, '{file_9.bin}', '{}', 'CyberKid', '{df,reversing}', '{}', 'http', 80, ''),
(10, 'challenge_10', 2, 'Solve challenge_10!', 'isolet-dev{random_flag_10}', 'dynamic', 300, '{file_10.bin}', '{}', 'Naughtyb0y', '{df,pwn}', '{}', 'nc', 30840, ''),
(11, 'challenge_11', 5, 'Solve challenge_11!', 'isolet-dev{random_flag_11}', 'static', 100, '{file_11.bin}', '{}', 'TheAlpha', '{df}', '{}', 'http', 80, ''),
(12, 'challenge_12', 3, 'Solve challenge_12!', 'isolet-dev{random_flag_12}', 'on-demand', 100, '{}', '{}', 'Naughtyb0y', '{linux}', '{}', 'http', 443, 'subd_12'),
(13, 'challenge_13', 8, 'Solve challenge_13!', 'isolet-dev{random_flag_13}', 'on-demand', 300, '{}', '{}', 'CyberKid', '{crypto,pwn,forensics}', '{}', 'http', 443, 'subd_13'),
(14, 'challenge_14', 9, 'Solve challenge_14!', 'isolet-dev{random_flag_14}', 'static', 100, '{file_14.bin}', '{}', 'TheAlpha', '{df,linux}', '{}', 'http', 80, ''),
(15, 'challenge_15', 5, 'Solve challenge_15!', 'isolet-dev{random_flag_15}', 'static', 100, '{file_15.bin}', '{}', 'TheAlpha', '{forensics,heap}', '{}', 'http', 80, ''),
(16, 'challenge_16', 8, 'Solve challenge_16!', 'isolet-dev{random_flag_16}', 'static', 100, '{file_16.bin}', '{}', 'TheAlpha', '{web,df}', '{}', 'http', 80, ''),
(17, 'challenge_17', 8, 'Solve challenge_17!', 'isolet-dev{random_flag_17}', 'static', 200, '{file_17.bin}', '{}', 'CyberKid', '{pwn,df,web}', '{}', 'http', 80, ''),
(18, 'challenge_18', 7, 'Solve challenge_18!', 'isolet-dev{random_flag_18}', 'on-demand', 200, '{}', '{}', 'TheAlpha', '{python,crypto,pwn}', '{}', 'http', 443, 'subd_18'),
(19, 'challenge_19', 7, 'Solve challenge_19!', 'isolet-dev{random_flag_19}', 'dynamic', 100, '{file_19.bin}', '{}', 'e4stw1nd', '{web}', '{}', 'nc', 30896, ''),
(20, 'challenge_20', 9, 'Solve challenge_20!', 'isolet-dev{random_flag_20}', 'static', 500, '{file_20.bin}', '{}', 'e4stw1nd', '{forensics}', '{}', 'http', 80, ''),
(21, 'challenge_21', 4, 'Solve challenge_21!', 'isolet-dev{random_flag_21}', 'dynamic', 300, '{file_21.bin}', '{}', 'CyberKid', '{heap}', '{}', 'http', 80, ''),
(22, 'challenge_22', 4, 'Solve challenge_22!', 'isolet-dev{random_flag_22}', 'static', 300, '{file_22.bin}', '{}', 'TheAlpha', '{misc,reversing}', '{}', 'http', 80, ''),
(23, 'challenge_23', 9, 'Solve challenge_23!', 'isolet-dev{random_flag_23}', 'dynamic', 100, '{file_23.bin}', '{}', 'e4stw1nd', '{linux,reversing,heap}', '{}', 'http', 80, ''),
(24, 'challenge_24', 2, 'Solve challenge_24!', 'isolet-dev{random_flag_24}', 'dynamic', 100, '{file_24.bin}', '{}', 'Hackerman', '{crypto,reversing,web}', '{}', 'http', 80, ''),
(25, 'challenge_25', 2, 'Solve challenge_25!', 'isolet-dev{random_flag_25}', 'static', 200, '{file_25.bin}', '{}', 'e4stw1nd', '{forensics,reversing}', '{}', 'http', 80, ''),
(26, 'challenge_26', 9, 'Solve challenge_26!', 'isolet-dev{random_flag_26}', 'on-demand', 300, '{}', '{}', 'e4stw1nd', '{python,forensics,web}', '{}', 'http', 443, 'subd_26'),
(27, 'challenge_27', 7, 'Solve challenge_27!', 'isolet-dev{random_flag_27}', 'on-demand', 400, '{}', '{}', 'Naughtyb0y', '{linux}', '{}', 'http', 443, 'subd_27'),
(28, 'challenge_28', 8, 'Solve challenge_28!', 'isolet-dev{random_flag_28}', 'on-demand', 200, '{}', '{}', 'Naughtyb0y', '{misc,web}', '{}', 'http', 443, 'subd_28'),
(29, 'challenge_29', 3, 'Solve challenge_29!', 'isolet-dev{random_flag_29}', 'on-demand', 300, '{}', '{}', 'Hackerman', '{python,reversing}', '{}', 'http', 443, 'subd_29'),
(30, 'challenge_30', 8, 'Solve challenge_30!', 'isolet-dev{random_flag_30}', 'static', 400, '{file_30.bin}', '{}', 'CyberKid', '{python}', '{}', 'http', 80, ''),
(31, 'challenge_31', 6, 'Solve challenge_31!', 'isolet-dev{random_flag_31}', 'dynamic', 300, '{file_31.bin}', '{}', 'TheAlpha', '{misc,df,crypto}', '{}', 'http', 80, ''),
(32, 'challenge_32', 8, 'Solve challenge_32!', 'isolet-dev{random_flag_32}', 'static', 400, '{file_32.bin}', '{}', 'TheAlpha', '{crypto,web}', '{}', 'http', 80, ''),
(33, 'challenge_33', 8, 'Solve challenge_33!', 'isolet-dev{random_flag_33}', 'static', 300, '{file_33.bin}', '{}', 'TheAlpha', '{forensics,crypto,pwn}', '{}', 'http', 80, ''),
(34, 'challenge_34', 8, 'Solve challenge_34!', 'isolet-dev{random_flag_34}', 'on-demand', 400, '{}', '{}', 'CyberKid', '{misc,df,python}', '{}', 'http', 443, 'subd_34'),
(35, 'challenge_35', 1, 'Solve challenge_35!', 'isolet-dev{random_flag_35}', 'static', 300, '{file_35.bin}', '{}', 'Naughtyb0y', '{forensics}', '{}', 'http', 80, ''),
(36, 'challenge_36', 1, 'Solve challenge_36!', 'isolet-dev{random_flag_36}', 'static', 500, '{file_36.bin}', '{}', 'e4stw1nd', '{heap,linux}', '{}', 'http', 80, ''),
(37, 'challenge_37', 9, 'Solve challenge_37!', 'isolet-dev{random_flag_37}', 'dynamic', 100, '{file_37.bin}', '{}', 'e4stw1nd', '{python,forensics,pwn}', '{}', 'nc', 30900, ''),
(38, 'challenge_38', 9, 'Solve challenge_38!', 'isolet-dev{random_flag_38}', 'on-demand', 200, '{}', '{}', 'e4stw1nd', '{df}', '{}', 'http', 443, 'subd_38'),
(39, 'challenge_39', 7, 'Solve challenge_39!', 'isolet-dev{random_flag_39}', 'static', 200, '{file_39.bin}', '{}', 'TheAlpha', '{web,python}', '{}', 'http', 80, ''),
(40, 'challenge_40', 9, 'Solve challenge_40!', 'isolet-dev{random_flag_40}', 'static', 400, '{file_40.bin}', '{}', 'Naughtyb0y', '{linux,crypto}', '{}', 'http', 80, ''),
(41, 'challenge_41', 9, 'Solve challenge_41!', 'isolet-dev{random_flag_41}', 'static', 200, '{file_41.bin}', '{}', 'Hackerman', '{df,web}', '{}', 'http', 80, ''),
(42, 'challenge_42', 2, 'Solve challenge_42!', 'isolet-dev{random_flag_42}', 'static', 500, '{file_42.bin}', '{}', 'CyberKid', '{pwn,reversing}', '{}', 'http', 80, ''),
(43, 'challenge_43', 8, 'Solve challenge_43!', 'isolet-dev{random_flag_43}', 'static', 200, '{file_43.bin}', '{}', 'CyberKid', '{heap,web}', '{}', 'http', 80, ''),
(44, 'challenge_44', 4, 'Solve challenge_44!', 'isolet-dev{random_flag_44}', 'on-demand', 400, '{}', '{}', 'e4stw1nd', '{misc,pwn}', '{}', 'http', 443, 'subd_44'),
(45, 'challenge_45', 2, 'Solve challenge_45!', 'isolet-dev{random_flag_45}', 'on-demand', 300, '{}', '{}', 'CyberKid', '{web,misc}', '{}', 'http', 443, 'subd_45'),
(46, 'challenge_46', 4, 'Solve challenge_46!', 'isolet-dev{random_flag_46}', 'static', 100, '{file_46.bin}', '{}', 'Hackerman', '{forensics}', '{}', 'http', 80, ''),
(47, 'challenge_47', 2, 'Solve challenge_47!', 'isolet-dev{random_flag_47}', 'static', 100, '{file_47.bin}', '{}', 'e4stw1nd', '{crypto,df}', '{}', 'http', 80, ''),
(48, 'challenge_48', 6, 'Solve challenge_48!', 'isolet-dev{random_flag_48}', 'static', 400, '{file_48.bin}', '{}', 'e4stw1nd', '{misc}', '{}', 'http', 80, ''),
(49, 'challenge_49', 2, 'Solve challenge_49!', 'isolet-dev{random_flag_49}', 'static', 400, '{file_49.bin}', '{}', 'CyberKid', '{misc,crypto}', '{}', 'http', 80, ''),
(50, 'challenge_50', 2, 'Solve challenge_50!', 'isolet-dev{random_flag_50}', 'on-demand', 400, '{}', '{}', 'Naughtyb0y', '{df}', '{}', 'http', 443, 'subd_50'),
(51, 'challenge_51', 9, 'Solve challenge_51!', 'isolet-dev{random_flag_51}', 'static', 500, '{file_51.bin}', '{}', 'TheAlpha', '{misc}', '{}', 'http', 80, ''),
(52, 'challenge_52', 7, 'Solve challenge_52!', 'isolet-dev{random_flag_52}', 'dynamic', 200, '{file_52.bin}', '{}', 'TheAlpha', '{linux}', '{}', 'nc', 30916, ''),
(53, 'challenge_53', 3, 'Solve challenge_53!', 'isolet-dev{random_flag_53}', 'dynamic', 300, '{file_53.bin}', '{}', 'CyberKid', '{python}', '{}', 'http', 80, ''),
(54, 'challenge_54', 7, 'Solve challenge_54!', 'isolet-dev{random_flag_54}', 'static', 300, '{file_54.bin}', '{}', 'Hackerman', '{misc,forensics}', '{}', 'http', 80, ''),
(55, 'challenge_55', 5, 'Solve challenge_55!', 'isolet-dev{random_flag_55}', 'static', 500, '{file_55.bin}', '{}', 'e4stw1nd', '{reversing}', '{}', 'http', 80, ''),
(56, 'challenge_56', 9, 'Solve challenge_56!', 'isolet-dev{random_flag_56}', 'dynamic', 500, '{file_56.bin}', '{}', 'Hackerman', '{linux}', '{}', 'nc', 30867, ''),
(57, 'challenge_57', 4, 'Solve challenge_57!', 'isolet-dev{random_flag_57}', 'on-demand', 400, '{}', '{}', 'CyberKid', '{crypto,heap}', '{}', 'http', 443, 'subd_57'),
(58, 'challenge_58', 1, 'Solve challenge_58!', 'isolet-dev{random_flag_58}', 'static', 100, '{file_58.bin}', '{}', 'e4stw1nd', '{reversing}', '{}', 'http', 80, ''),
(59, 'challenge_59', 6, 'Solve challenge_59!', 'isolet-dev{random_flag_59}', 'dynamic', 200, '{file_59.bin}', '{}', 'CyberKid', '{reversing,df}', '{}', 'http', 80, ''),
(60, 'challenge_60', 8, 'Solve challenge_60!', 'isolet-dev{random_flag_60}', 'dynamic', 200, '{file_60.bin}', '{}', 'e4stw1nd', '{web}', '{}', 'http', 80, '');

-- Add hints
INSERT INTO hints (hid, chall_id, hint, cost) VALUES
(1, 2, 'Hint for challenge_2', 10),
(2, 2, 'Hint for challenge_2', 20),
(3, 3, 'Hint for challenge_3', 20),
(4, 3, 'Hint for challenge_3', 20),
(5, 3, 'Hint for challenge_3', 20),
(6, 6, 'Hint for challenge_6', 30),
(7, 6, 'Hint for challenge_6', 30),
(8, 8, 'Hint for challenge_8', 20),
(9, 8, 'Hint for challenge_8', 10),
(10, 9, 'Hint for challenge_9', 10),
(11, 10, 'Hint for challenge_10', 10),
(12, 11, 'Hint for challenge_11', 10),
(13, 11, 'Hint for challenge_11', 20),
(14, 14, 'Hint for challenge_14', 20),
(15, 15, 'Hint for challenge_15', 30),
(16, 16, 'Hint for challenge_16', 30),
(17, 16, 'Hint for challenge_16', 30),
(18, 17, 'Hint for challenge_17', 30),
(19, 18, 'Hint for challenge_18', 10),
(20, 20, 'Hint for challenge_20', 30),
(21, 20, 'Hint for challenge_20', 20),
(22, 21, 'Hint for challenge_21', 30),
(23, 22, 'Hint for challenge_22', 30),
(24, 22, 'Hint for challenge_22', 10),
(25, 22, 'Hint for challenge_22', 10),
(26, 23, 'Hint for challenge_23', 30),
(27, 24, 'Hint for challenge_24', 30),
(28, 24, 'Hint for challenge_24', 10),
(29, 26, 'Hint for challenge_26', 30),
(30, 27, 'Hint for challenge_27', 20),
(31, 28, 'Hint for challenge_28', 30),
(32, 28, 'Hint for challenge_28', 30),
(33, 30, 'Hint for challenge_30', 30),
(34, 31, 'Hint for challenge_31', 10),
(35, 32, 'Hint for challenge_32', 30),
(36, 33, 'Hint for challenge_33', 30),
(37, 35, 'Hint for challenge_35', 20),
(38, 35, 'Hint for challenge_35', 30),
(39, 36, 'Hint for challenge_36', 30),
(40, 36, 'Hint for challenge_36', 20),
(41, 36, 'Hint for challenge_36', 10),
(42, 40, 'Hint for challenge_40', 10),
(43, 40, 'Hint for challenge_40', 20),
(44, 41, 'Hint for challenge_41', 30),
(45, 43, 'Hint for challenge_43', 10),
(46, 43, 'Hint for challenge_43', 30),
(47, 44, 'Hint for challenge_44', 10),
(48, 44, 'Hint for challenge_44', 20),
(49, 48, 'Hint for challenge_48', 30),
(50, 48, 'Hint for challenge_48', 20),
(51, 49, 'Hint for challenge_49', 10),
(52, 50, 'Hint for challenge_50', 20),
(53, 51, 'Hint for challenge_51', 20),
(54, 52, 'Hint for challenge_52', 30),
(55, 55, 'Hint for challenge_55', 10),
(56, 56, 'Hint for challenge_56', 10),
(57, 57, 'Hint for challenge_57', 10),
(58, 57, 'Hint for challenge_57', 30),
(59, 60, 'Hint for challenge_60', 30),
(60, 60, 'Hint for challenge_60', 20);

INSERT INTO users (userid, email, username, password, teamid) VALUES 
(1, 'menoexist@gmail.com', 'standalone',  'b00cf059816d1d134ba722b08b3e330cd9c229ff8faa07a40ed4c795917a23a4', 1),
(2, 'live.stream@gmail.com', 'getone', 'b00cf059816d1d134ba722b08b3e330cd9c229ff8faa07a40ed4c795917a23a4', 2),
(3, 'kingkongvsgodzilla@gmail.com', 'uselesscop', 'b00cf059816d1d134ba722b08b3e330cd9c229ff8faa07a40ed4c795917a23a4', 1),
(4, 'lord.shardul@gmail.com', 'rayofhope', 'b00cf059816d1d134ba722b08b3e330cd9c229ff8faa07a40ed4c795917a23a4', 1),
(5, 'getthisout@gmail.com', 'glimpse', 'b00cf059816d1d134ba722b08b3e330cd9c229ff8faa07a40ed4c795917a23a4', 2),
(6, 'infoseccl@gmail.com', 'cyberlabs', 'b00cf059816d1d134ba722b08b3e330cd9c229ff8faa07a40ed4c795917a23a4', 2),
(7, 'admin@isolet-dev.in', 'admin', 'b00cf059816d1d134ba722b08b3e330cd9c229ff8faa07a40ed4c795917a23a4', 1);

-- Teams
INSERT INTO teams (teamid, teamname, captain, password) VALUES
(1, 'TitanCrew', 1, 'b00cf059816d1d134ba722b08b3e330cd9c229ff8faa07a40ed4c795917a23a4'),
(2, 'BIT CRIMINALS', 2, 'b00cf059816d1d134ba722b08b3e330cd9c229ff8faa07a40ed4c795917a23a4');

SELECT setval('users_userid_seq', 7);
SELECT setval('teams_teamid_seq', 2);

UPDATE challenges SET visible = true;
UPDATE hints SET visible = true;

SELECT setval('categories_category_id_seq', 4);
SELECT setval('challenges_chall_id_seq', 60);
SELECT setval('hints_hid_seq', 69);

INSERT INTO config (key, value) VALUES
('EVENT_START', '1728813375'),
('EVENT_END', '1799986175'),
('POST_EVENT', 'false'),
('TEAM_LEN', '4'),
('CTF_NAME', 'isolet-dev'),
('CONCURRENT_INSTANCES', '1'),
('INSTANCE_HOSTNAME', 'ctf.isolet-dev.in'),
('INSTANCE_TIME', '30'),
('MAX_INSTANCE_TIME', '60'),
('EMAIL_ID', 'menoone@isolet-dev.in'),
('EMAIL_AUTH', 'jdrglctcdtdlbi'),
('PUBLIC_URL', 'isolet-dev.in'),
('SMTP_HOST', 'smtp.gmail.com'),
('SMTP_PORT', '587'),
('SYNC_CONFIG_SECONDS', '30'),
('USER_REGISTRATION', 'true'),
('API_RATE_LIMIT', 'false'),
('EMAIL_VERIFICATION', 'false'),
('EMAIL_USERNAME', 'something');