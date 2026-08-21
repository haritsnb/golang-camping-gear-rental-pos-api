-- +migrate Up

INSERT INTO roles (id, name, description, created_at, modified_at)
VALUES 
    (1, 'admin', 'Super Administrator dengan hak akses penuh ke seluruh modul sistem', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, 'cashier', 'Kasir untuk operasional transaksi booking, handover, return, dan pembayaran', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (3, 'staff', 'Staff operasional gudang, inventaris produk, dan pemeliharaan alat', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (name) DO NOTHING;

SELECT setval('roles_id_seq', (SELECT MAX(id) FROM roles));

INSERT INTO categories (id, name, description, created_at, modified_at)
VALUES 
    (1, 'Tenda & Shelter', 'Tenda dome, flysheet, tarp shelter, dan perlengkapan pasak', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, 'Sleeping Gear', 'Sleeping bag, matras tiup/foil, hammock, dan bantal angin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (3, 'Cooking & Dining', 'Kompor portable, nesting, cooking set, pisau lipat, dan gas canister', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (4, 'Carrier & Ransel', 'Carrier ransel gunung 40L - 80L, rain cover, dan daypack', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (5, 'Lighting & Electronics', 'Headlamp, lentera tenda, senter LED, dan power station outdoor', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (6, 'Trekking & Safety', 'Trekking pole, survival kit, peluit, jas hujan, dan first aid kit', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (name) DO NOTHING;

SELECT setval('categories_id_seq', (SELECT MAX(id) FROM categories));

INSERT INTO brands (id, name, description, is_active, created_at, modified_at)
VALUES 
    (1, 'Eiger', 'Brand perlengkapan outdoor lokal legendaris dan premium', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, 'Naturehike', 'Brand spesialis ultralight camping gear dan tenda modern', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (3, 'Consina', 'Brand outdoor gear populer, tangguh, dan terjangkau', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (4, 'Arei Outdoor Gear', 'Produsen perlengkapan petualang dan pendakian gunung', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (5, 'Deuter', 'Brand asal Jerman spesialis tas carrier ransel ergonomis kelas dunia', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (6, 'Kovea', 'Brand internasional spesialis kompor portable dan burner outdoor', TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (name) DO NOTHING;

SELECT setval('brands_id_seq', (SELECT MAX(id) FROM brands));

INSERT INTO products (id, category_id, brand_id, name, rental_price_per_day, default_deposit, late_fee_per_day, lost_compensation_fee, is_active, created_at, modified_at)
VALUES 
    (1, 1, 1, 'Tenda Eiger Guardian 4P (Kapasitas 4 Orang)', 45000.00, 100000.00, 25000.00, 650000.00, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, 1, 2, 'Tenda Naturehike Cloud Up 2 Ultralight (2P)', 40000.00, 100000.00, 20000.00, 750000.00, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (3, 2, 3, 'Sleeping Bag Consina Mummy Snuggle 300', 15000.00, 30000.00, 10000.00, 180000.00, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (4, 3, 6, 'Kompor Portable Kovea Spider Burner', 20000.00, 50000.00, 15000.00, 320000.00, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (5, 3, 4, 'Cooking Set Nesting Arei DS-308 (4 in 1)', 15000.00, 30000.00, 10000.00, 210000.00, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (6, 4, 5, 'Carrier Deuter Aircontact Core 60+10L', 50000.00, 150000.00, 30000.00, 1800000.00, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (7, 5, 2, 'Headlamp Naturehike LED 200 Lumens Waterproof', 10000.00, 20000.00, 5000.00, 120000.00, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (8, 6, 1, 'Trekking Pole Eiger Carbon Shock Absorber', 15000.00, 30000.00, 10000.00, 250000.00, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

SELECT setval('products_id_seq', (SELECT MAX(id) FROM products));

INSERT INTO product_units (id, product_id, unit_code, serial_number, condition, status, notes, created_at, modified_at)
VALUES 
    (1, 1, 'TND-EIG-001', 'SN-EIG-2026-001', 'good', 'available', 'Kondisi mulus, pasak lengkap 10 pcs, frame alloy utuh', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, 1, 'TND-EIG-002', 'SN-EIG-2026-002', 'good', 'available', 'Kondisi mulus, lengkap dengan footprint & guyline', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (3, 2, 'TND-NH-001', 'SN-NH-2026-088', 'good', 'available', 'Tenda ultralight warna hijau army, pasak aluminium 11 pcs', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (4, 3, 'SB-CSN-001', 'SN-CSN-SB-101', 'good', 'available', 'Sudah dilaundry bersih, wangi, resleting lancar', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (5, 3, 'SB-CSN-002', 'SN-CSN-SB-102', 'good', 'available', 'Sudah dilaundry bersih, wangi, resleting lancar', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (6, 4, 'KMP-KOV-001', 'SN-KOV-SP-01', 'good', 'available', 'Pemantik elektrik berfungsi normal, selang fleksibel bersih', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (7, 5, 'NST-ARI-001', 'SN-ARI-DS-01', 'good', 'available', 'Lengkap 3 panci + 1 wajan + sponge cuci + tas jaring', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (8, 6, 'CRR-DTR-001', 'SN-DTR-AC-001', 'good', 'available', 'Include original Rain Cover warna biru, busa backsystem tebal', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (9, 7, 'HL-NH-001', 'SN-NH-HL-301', 'good', 'available', 'Baterai rechargeable terisi 100%, strap elastis bagus', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (10, 8, 'TP-EIG-001', 'SN-EIG-TP-99', 'good', 'available', 'Sepasang (2 pcs), mekanisme antishock & lock system normal', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (unit_code) DO NOTHING;

SELECT setval('product_units_id_seq', (SELECT MAX(id) FROM product_units));

INSERT INTO customers (id, name, identity_type, identity_number, identity_photo_url, phone, emergency_contact, email, address, is_blacklisted, notes, created_at, modified_at)
VALUES 
    (1, 'Budi Pratama', 'KTP', '3201012304950001', '/storages/identities/sample_ktp_budi.jpg', '081233445566', '081299887766 (Istri)', 'budi.pratama@example.com', 'Jl. Sukajadi No. 128, Sukasari, Kota Bandung, Jawa Barat', FALSE, 'Pelanggan tetap komunitas pendaki Bandung', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (2, 'Siti Nurhaliza', 'SIM', '980712345678', '/storages/identities/sample_sim_siti.jpg', '085711223344', '085799001122 (Kakak)', 'siti.nurhaliza@example.com', 'Jl. Kaliurang KM 5.5, Depok, Sleman, D.I. Yogyakarta', FALSE, 'Verifikasi identitas SIM A valid', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (identity_number) DO NOTHING;

SELECT setval('customers_id_seq', (SELECT MAX(id) FROM customers));


-- +migrate Down

DELETE FROM customers WHERE identity_number IN ('3201012304950001', '980712345678');
DELETE FROM product_units WHERE id BETWEEN 1 AND 10;
DELETE FROM products WHERE id BETWEEN 1 AND 8;
DELETE FROM brands WHERE id BETWEEN 1 AND 6;
DELETE FROM categories WHERE id BETWEEN 1 AND 6;
DELETE FROM roles WHERE name IN ('admin', 'cashier', 'staff');