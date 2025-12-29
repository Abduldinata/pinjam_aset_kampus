--
-- PostgreSQL database dump (Cleaned & Portable)
-- Maintenance oleh: Abduldinata (Tim 3)
-- Forked from: Yeflou
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', 'public', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

-- 1. Create Types (Enum-like) with Safety Checks
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role_type') THEN
        CREATE TYPE public.role_type AS ENUM ('admin', 'user');
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status_type') THEN
        CREATE TYPE public.status_type AS ENUM ('dipinjam', 'kembali', 'terlambat');
    END IF;
END $$;

SET default_tablespace = '';
SET default_table_access_method = heap;

-- 2. Create Tables with IF NOT EXISTS
CREATE TABLE IF NOT EXISTS public.items (
    id SERIAL PRIMARY KEY,
    name character varying(100),
    category character varying(50),
    stock integer DEFAULT 0 NOT NULL,
    location character varying(100),
    description text,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

CREATE TABLE IF NOT EXISTS public.users (
    id SERIAL PRIMARY KEY,
    name character varying(100),
    email character varying(100) UNIQUE,
    password character varying(255),
    role public.role_type DEFAULT 'user'::public.role_type,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

CREATE TABLE IF NOT EXISTS public.loans (
    id SERIAL PRIMARY KEY,
    user_id integer REFERENCES public.users(id) ON DELETE CASCADE,
    item_id integer REFERENCES public.items(id) ON DELETE RESTRICT,
    borrow_date date,
    due_date date,
    return_date date,
    status public.status_type DEFAULT 'dipinjam'::public.status_type,
    notes text,
    is_fine_paid boolean DEFAULT false,
    payment_method character varying(50),
    payment_proof text,
    fine_amount bigint DEFAULT 0,
    student_id_card text,
    created_at timestamp without time zone,
    updated_at timestamp without time zone
);

CREATE TABLE IF NOT EXISTS public.notifications (
    id SERIAL PRIMARY KEY,
    user_id integer REFERENCES public.users(id) ON DELETE CASCADE,
    loan_id integer REFERENCES public.loans(id) ON DELETE SET NULL,
    message text,
    is_read boolean DEFAULT false,
    type character varying(50) DEFAULT 'reminder'::character varying,
    created_at timestamp without time zone
);

-- 3. Data for Name: items
INSERT INTO public.items (id, name, category, stock, location, description, created_at, updated_at) VALUES
(7, 'Kunci Lab IOT', 'Lainnya', 1, 'Front Office', 'sdsdf', '2025-12-24 07:48:26', '2025-12-24 07:48:26'),
(5, 'Laptop Asus New', 'Elektronik', 4, 'Lab Komputer', 'Laptop Praktikum', '2025-12-24 00:35:46', '2025-12-24 17:16:59'),
(8, 'Buku Jurnal Pencatatan', 'ATK', 3, 'Ruang Admin', 'Buku Jurnal Pencatatan', '2025-12-24 17:43:46', '2025-12-24 17:43:46'),
(9, 'Kabel Charge Type C', 'Elektronik', 5, 'Ruang Teknisi', 'Kabel Charge Type C', '2025-12-24 17:45:08', '2025-12-24 17:45:08'),
(2, 'Kabel HDMI', 'Aksesoris', 9, 'Ruang TU', 'Panjang 5 meter', '2025-12-23 23:23:47', '2025-12-29 02:37:46'),
(1, 'Proyektor Epson', 'Elektronik', 5, 'Ruang TU', 'Proyektor HD', '2025-12-23 23:23:47', '2025-12-29 03:31:09'),
(6, 'Kunci Lab 5', 'Lainnya', 0, 'Ruang Admin', 'ini kunci', '2025-12-24 01:03:45', '2025-12-29 04:48:30'),
(3, 'Sound Portable', 'Audio', 2, 'Gudang', 'Speaker untuk acara', '2025-12-23 23:23:47', '2025-12-29 06:10:21')
ON CONFLICT (id) DO UPDATE SET
name = EXCLUDED.name, category = EXCLUDED.category, stock = EXCLUDED.stock,
location = EXCLUDED.location, description = EXCLUDED.description, updated_at = EXCLUDED.updated_at;

-- 4. Data for Name: users
-- Password default: 'admin123' & 'mhs123' (sudah di hash bcrypt)
INSERT INTO public.users (id, name, email, password, role, created_at, updated_at) VALUES
(3, 'Super Admin', 'super@admin.com', '$2a$10$7/5Fdu2kcGwEZobcXb.o.eQAGqyG/DJ0tmKvZbY2J.VHAodtKB0N.', 'admin', '2025-12-24 00:02:41', '2025-12-24 00:02:41'),
(4, 'Amelia', 'maba@kampus.ac.id', '$2a$10$LRUzY9IlSupQgzfJw4rgYOvoxkcLui0/hXgrhO8BO8bT.gg1GlXs6', 'user', '2025-12-24 00:50:42', '2025-12-24 00:50:42'),
(5, 'Test User', 'testuser@example.com', '$2a$10$d1ngkycBmG8ul7D5/XLLNOQgjMVGtGdnILiEq7YJ.bDQvN8X.1Usy', 'user', '2025-12-24 07:54:37', '2025-12-24 07:54:37'),
(1, 'Admin Kampus', 'admin@kampus.ac.id', '$2a$10$Ijdvmvp/VPCsQ.Utp0c87u9FO5bv9aKNmSPMull2K7V.ZSzDBcsYq', 'admin', '2025-12-23 23:23:47', '2025-12-23 23:23:47'),
(2, 'Budi Mahasiswa', 'budi@mhs.ac.id', '$2a$10$1D.2XWzIrh0PaaCqqkK.NOoocXEKZUZOQkvG7aSUcnPa8E.UNyVbW', 'user', '2025-12-23 23:23:47', '2025-12-23 23:23:47')
ON CONFLICT (id) DO UPDATE SET
name = EXCLUDED.name, email = EXCLUDED.email, password = EXCLUDED.password, role = EXCLUDED.role, updated_at = EXCLUDED.updated_at;

-- 5. Data for Name: loans
INSERT INTO public.loans (id, user_id, item_id, borrow_date, due_date, return_date, status, notes, created_at, updated_at, is_fine_paid, payment_method, payment_proof, fine_amount, student_id_card) VALUES
(2, 4, 1, '2025-12-23', '2025-12-24', NULL, 'dipinjam', 'Presentasi Tugas', '2025-12-24 03:57:39', '2025-12-24 03:57:39', false, NULL, NULL, 0, NULL),
(3, 4, 3, '2025-12-23', '2025-12-25', '2025-12-23', 'kembali', 'pensi', '2025-12-24 05:18:50', '2025-12-24 05:35:05', false, NULL, NULL, 0, NULL),
(4, 4, 5, '2025-12-24', '2025-12-25', NULL, 'dipinjam', 'tugas akhir', '2025-12-24 17:16:59', '2025-12-24 17:16:59', false, NULL, NULL, 0, NULL),
(5, 2, 2, '2025-12-28', '2025-12-29', NULL, 'dipinjam', 'tugas golang', '2025-12-29 02:37:46', '2025-12-29 02:37:46', false, NULL, NULL, 0, NULL),
(1, 2, 1, '2023-12-01', '2023-12-03', '2025-12-28', 'kembali', ' [Dikembalikan Terlambat]', '2025-12-23 23:23:47', '2025-12-29 03:31:09', false, NULL, NULL, 0, NULL),
(6, 2, 6, '2025-12-28', '2025-12-29', NULL, 'dipinjam', 'pinjam', '2025-12-29 04:48:30', '2025-12-29 04:48:30', false, '', '', 0, NULL),
(7, 2, 3, '2025-12-28', '2025-12-29', '2025-12-28', 'kembali', 'tugas akhir', '2025-12-29 06:07:35', '2025-12-29 06:10:21', false, '', '', 0, 'user_2/ktm_1766938055_Screenshot (4).png')
ON CONFLICT (id) DO UPDATE SET
user_id = EXCLUDED.user_id, item_id = EXCLUDED.item_id, status = EXCLUDED.status, notes = EXCLUDED.notes, updated_at = EXCLUDED.updated_at;

-- 6. Data for Name: notifications
INSERT INTO public.notifications (id, user_id, message, is_read, created_at, loan_id, type) VALUES
(1, 4, 'Halo Amelia, mohon segera kembalikan aset: Proyektor Epson. Batas waktu: 24 Dec 2025', false, '2025-12-24 06:11:47', NULL, 'reminder'),
(2, 4, 'Halo Amelia, mohon segera kembalikan aset: Laptop Asus New. Batas waktu: 25 Dec 2025', false, '2025-12-24 17:17:34', NULL, 'reminder'),
(5, 2, 'Halo Budi Mahasiswa, mohon segera kembalikan aset: Kabel HDMI. Batas waktu: 29 Dec 2025', true, '2025-12-29 03:29:52', 0, 'reminder'),
(4, 2, 'Halo Budi Mahasiswa, mohon segera kembalikan aset: Kabel HDMI. Batas waktu: 29 Dec 2025', true, '2025-12-29 03:29:31', 0, 'reminder'),
(3, 2, '⚠️ TERLAMBAT! Aset [Proyektor Epson] belum dikembalikan. Batas waktu: 03 Dec 2023. Segera kembalikan ke Admin!', true, '2025-12-29 02:37:34', 1, 'denda')
ON CONFLICT (id) DO UPDATE SET
message = EXCLUDED.message, is_read = EXCLUDED.is_read;

-- 7. Reset Sequences
SELECT pg_catalog.setval('public.items_id_seq', (SELECT MAX(id) FROM public.items), true);
SELECT pg_catalog.setval('public.loans_id_seq', (SELECT MAX(id) FROM public.loans), true);
SELECT pg_catalog.setval('public.notifications_id_seq', (SELECT MAX(id) FROM public.notifications), true);
SELECT pg_catalog.setval('public.users_id_seq', (SELECT MAX(id) FROM public.users), true);
db_peminjaman_kampus