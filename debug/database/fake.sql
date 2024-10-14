-- Create categories
INSERT INTO categories (category_id, category_name) VALUES 
    (1, 'Miscellaneous'),
    (2, 'Reversing'),
    (3, 'Pwn'),
    (4, 'Linux');

-- Add challenge data
INSERT INTO challenges (
    chall_id, chall_name, category_id, 
    prompt, flag, type, 
    points, files, requirements, 
    author, tags, links
) VALUES 
(
    1, 'byteme', 2,
    'I know you are a python expert, but can you reverse this?',
    'isolet-dev{e4sy_p34sy_byt3c0d3_d1s4sm}',
    'static',
    300, '{byteme.pyc}', '{}',
    'TheAlpha', '{python,bytecode}', '{}'
),
(
    2, 'babyheap', 3,
    'Just a normal note taking app...',
    'isolet-dev{He4p_1s_5ecur3_f0r_n0t3s??}',
    'dynamic',
    400, '{heap,libc.so.6}', '{}',
    'Naughtyb0y', '{heap,pwn}', '{https://ctf101.org/binary-exploitation/heap-exploitation/}'
),
(
    3, 'flag-finder', 3,
    'I am wondering can you find the needle in the haystack (Not with your eyes but just binary)???',
    'isolet-dev{f1nd_f1nd_f1nding_th3_fl4ggggg!}',
    'dynamic',
    400, '{flag-finder}', '{2}',
    'Naughtyb0y', '{pwn}', '{}'
),
(
    4, 'Machine Trouble', 1,
    'I dont think even smart people like you can get the flag out of such a machine with so little memory.',
    'isolet-dev{dfa_hacked}',
    'on-demand',
    250, '{}', '{}',
    'e4stw1nd', '{dfa,misc}', '{}'
),
(
    5, 'Unixit', 4,
    'Nothing..just to teach you how to ssh',
    'isolet-dev{n0w_y0u_kn0w}',
    'on-demand',
    100, '{}', '{}',
    'bitc', '{ssh,linux}', '{}'
),
(
    6, 'pwn-me', 4,
    'Lets see if you can exploit me',
    'isolet-dev{1_h4v3nt_exp3ct3d_th1s}',
    'on-demand',
    1000, '{}', '{5}',
    'anonymous', '{0day,linux}', '{}'
);

UPDATE challenges SET visible = true;

INSERT INTO hints (hid, chall_id, hint, cost) VALUES
    (1, 4, 'At a time the only information you get is if a string is a part of the language you defined', 30),
    (2, 6, 'Can never do this before event ends lol', 50),
    (3, 5, 'Lets see a free hint for you', 0);

UPDATE hints SET visible = true;

INSERT INTO images (chall_id, registry, image, deployment, port, subd, cpu, mem) VALUES 
(2, '', 'babyheap', 'nc', 31000, 'babyheap', 30, 32),
(3, '', 'flag-finder', 'nc', 31001, 'flag-finder', 30, 32),
(4, '', 'dfa', 'nc', 5000, 'dfa', 50, 128),
(5, 'docker.io/infoseccl', 'unixit-level0', 'ssh', 22, 'unixit', 5, 10),
(6, 'public.ecr.aws/nginx/', 'nginx:alpine-slim', 'http', 80, 'nginx', 5, 10);

-- Users
INSERT INTO users (userid, email, username, password, teamid) VALUES 
(1, 'menoexist@gmail.com', 'standalone',  'b00cf059816d1d134ba722b08b3e330cd9c229ff8faa07a40ed4c795917a23a4', 1),
(2, 'live.stream@gmail.com', 'getone', 'b00cf059816d1d134ba722b08b3e330cd9c229ff8faa07a40ed4c795917a23a4', 1),
(3, 'kingkongvsgodzilla@gmail.com', 'uselesscop', 'b00cf059816d1d134ba722b08b3e330cd9c229ff8faa07a40ed4c795917a23a4', 1),
(4, 'lord.shardul@gmail.com', 'rayofhope', 'b00cf059816d1d134ba722b08b3e330cd9c229ff8faa07a40ed4c795917a23a4', 1),
(5, 'getthisout@gmail.com', 'glimpse', 'b00cf059816d1d134ba722b08b3e330cd9c229ff8faa07a40ed4c795917a23a4', 2),
(6, 'infoseccl@gmail.com', 'cyberlabs', 'b00cf059816d1d134ba722b08b3e330cd9c229ff8faa07a40ed4c795917a23a4', 2);

-- Teams
INSERT INTO teams (teamid, teamname, captain, members, password) VALUES
(1, 'TitanCrew', 1, '{2,3,4}', 'b00cf059816d1d134ba722b08b3e330cd9c229ff8faa07a40ed4c795917a23a4'),
(2, 'BIT CRIMINALS', 2, '{5,6}', 'b00cf059816d1d134ba722b08b3e330cd9c229ff8faa07a40ed4c795917a23a4');

-- Sublogs
INSERT INTO sublogs 
(chall_id, userid, teamid, flag, correct, ip) VALUES
(1, 2, 1, 'isolet-dev{this_is_not_correct}', false, '215.57.121.23'),
(1, 2, 1, 'isolet-dev{dont_brute_man}', false, '215.57.121.23'),
(1, 4, 1, 'isolet-dev{e4sy_p34sy_byt3c0d3_d1s4sm}', true, '151.64.234.12'),
(1, 5, 2, 'isolet-dev{e4sy_p34sy_byt3c0d3_d1s4sm}', true, '123.21.67.37'),
(1, 1, 1, 'isolet-dev{yep_not_a_flag}', false, '215.57.121.24'),
(2, 5, 2, 'isolet-dev{definitely_a_guess}', false, '123.21.67.37'),
(2, 6, 2, 'isolet-dev{He4p_1s_5ecur3_f0r_n0t3s??}', true, '123.21.67.38');

-- Hints
UPDATE teams SET uhints = '{1,2}' WHERE teamid = 1;
UPDATE teams SET uhints = '{2}' WHERE teamid = 2;