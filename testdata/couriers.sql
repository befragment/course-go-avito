-- Seed data for couriers
-- Phone format: +[0-9]{11}
-- transport_type: 'on_foot', 'scooter', 'car'

INSERT INTO couriers 
    (name,       phone,          status,      transport_type)
VALUES
    ('Vasya',    '+79000000001', 'available', 'scooter'),
    ('Petya',    '+79000000002', 'available', 'on_foot'),
    ('Masha',    '+79000000003', 'available', 'car'),
    ('Dasha',    '+79000000004', 'available', 'scooter'),
    ('Serega',   '+79000000005', 'available', 'on_foot'),
    ('Oleg',     '+79000000006', 'available', 'car'),
    ('Irina',    '+79000000007', 'available', 'scooter'),
    ('Nastya',   '+79000000008', 'available', 'on_foot'),
    ('Alex',     '+79000000009', 'available', 'car'),
    ('Nikolay',  '+79000000010', 'available', 'scooter'),
    ('Polina',   '+79000000011', 'available', 'on_foot'),
    ('Denis',    '+79000000012', 'available', 'car'),
    ('Sergey',   '+79000000013', 'available', 'scooter'),
    ('Alina',    '+79000000014', 'available', 'on_foot'),
    ('Roman',    '+79000000015', 'available', 'car'),
    ('Dmitry',   '+79000000016', 'available', 'scooter'),
    ('Olga',     '+79000000017', 'available', 'on_foot'),
    ('Yana',     '+79000000018', 'available', 'car'),
    ('Kirill',   '+79000000019', 'available', 'scooter'),
    ('Victoria', '+79000000020', 'available', 'on_foot'),
    ('Andrey',   '+79000000021', 'available', 'car'),
    ('Igor',     '+79000000022', 'available', 'scooter'),
    ('Elena',    '+79000000023', 'available', 'on_foot'),
    ('Tatiana',  '+79000000024', 'available', 'car');
