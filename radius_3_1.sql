-- phpMyAdmin SQL Dump
-- version 5.1.0
-- https://www.phpmyadmin.net/
--
-- Host: db_mysql_reseller
-- Waktu pembuatan: 26 Apr 2021 pada 14.46
-- Versi server: 8.0.24
-- Versi PHP: 7.4.16

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `radius`
--

-- --------------------------------------------------------

--
-- Struktur dari tabel `message`
--

CREATE TABLE `message` (
  `id` int NOT NULL,
  `chat_id` int NOT NULL,
  `received` enum('yes','no') NOT NULL DEFAULT 'no',
  `message` text NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `nas`
--

CREATE TABLE `nas` (
  `id` int NOT NULL,
  `nasname` varchar(128) NOT NULL,
  `shortname` varchar(32) DEFAULT NULL,
  `type` varchar(30) DEFAULT 'other',
  `ports` int DEFAULT NULL,
  `secret` varchar(60) NOT NULL DEFAULT 'secret',
  `server` varchar(64) DEFAULT NULL,
  `community` varchar(50) DEFAULT NULL,
  `description` varchar(200) DEFAULT 'RADIUS Client'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `package`
--

CREATE TABLE `package` (
  `id` int NOT NULL,
  `name` varchar(255) NOT NULL,
  `validity_value` int NOT NULL,
  `validity_unit` enum('HOUR','DAY','MONTH','YEAR') NOT NULL DEFAULT 'HOUR',
  `price` int NOT NULL,
  `margin` int NOT NULL,
  `profile` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `package`
--

INSERT INTO `package` (`id`, `name`, `validity_value`, `validity_unit`, `price`, `margin`, `profile`, `created_at`) VALUES
(1, 'paket 1 hari', 1, 'DAY', 2000, 1000, '2M_Profile', '2021-04-26 06:00:34');

-- --------------------------------------------------------

--
-- Struktur dari tabel `radacct`
--

CREATE TABLE `radacct` (
  `radacctid` bigint NOT NULL,
  `acctsessionid` varchar(64) NOT NULL DEFAULT '',
  `acctuniqueid` varchar(32) NOT NULL DEFAULT '',
  `username` varchar(64) NOT NULL DEFAULT '',
  `realm` varchar(64) DEFAULT '',
  `nasipaddress` varchar(15) NOT NULL DEFAULT '',
  `nasportid` varchar(15) DEFAULT NULL,
  `nasporttype` varchar(32) DEFAULT NULL,
  `acctstarttime` datetime DEFAULT NULL,
  `acctupdatetime` datetime DEFAULT NULL,
  `acctstoptime` datetime DEFAULT NULL,
  `acctinterval` int DEFAULT NULL,
  `acctsessiontime` int UNSIGNED DEFAULT NULL,
  `acctauthentic` varchar(32) DEFAULT NULL,
  `connectinfo_start` varchar(50) DEFAULT NULL,
  `connectinfo_stop` varchar(50) DEFAULT NULL,
  `acctinputoctets` bigint DEFAULT NULL,
  `acctoutputoctets` bigint DEFAULT NULL,
  `calledstationid` varchar(50) NOT NULL DEFAULT '',
  `callingstationid` varchar(50) NOT NULL DEFAULT '',
  `acctterminatecause` varchar(32) NOT NULL DEFAULT '',
  `servicetype` varchar(32) DEFAULT NULL,
  `framedprotocol` varchar(32) DEFAULT NULL,
  `framedipaddress` varchar(15) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Trigger `radacct`
--
DELIMITER $$
CREATE TRIGGER `bussiness_radacct_after_insert_trigger` AFTER INSERT ON `radacct` FOR EACH ROW BEGIN
		
	SET @expiration = (SELECT COUNT(*) FROM radcheck WHERE username = New.username AND attribute = 'Expiration'); 
		
	IF (@expiration = 0) THEN
		SET @validity_value = (SELECT package.validity_value FROM radpackage INNER JOIN package ON package_id = radpackage.package.id WHERE radpackage.username = New.username);
		SET @validity_unit = (SELECT package.validity_unit FROM radpackage INNER JOIN package ON package.id = radpackage.package_id WHERE radpackage.username = New.username);

		IF (@validity_unit = 'HOUR') THEN
			INSERT INTO radcheck(username, attribute, op, value) VALUES(New.username, "Expiration", ":=", DATE_FORMAT((NOW() + INTERVAL @validity_value HOUR), "%d %b %Y %H:%I:%S"));

		ELSEIF (@validity_unit = 'DAY') THEN
			INSERT INTO radcheck(username, attribute, op, value) VALUES(New.username, "Expiration", ":=", DATE_FORMAT((NOW() + INTERVAL @validity_value DAY), "%d %b %Y %H:%I:%S"));

		ELSEIF (@validity_unit = 'MONTH') THEN
			INSERT INTO radcheck(username, attribute, op, value) VALUES(New.username, "Expiration", ":=", DATE_FORMAT((NOW() + INTERVAL @validity_value MONTH), "%d %b %Y %H:%I:%S"));

		ELSEIF (@validity_unit = 'YEAR') THEN
			INSERT INTO radcheck(username, attribute, op, value) VALUES(New.username, "Expiration", ":=", DATE_FORMAT((NOW() + INTERVAL @validity_unit YEAR), "%d %b %Y %H:%I:%S"));

		END IF;

	END IF;
	END
$$
DELIMITER ;

-- --------------------------------------------------------

--
-- Struktur dari tabel `radcheck`
--

CREATE TABLE `radcheck` (
  `id` int UNSIGNED NOT NULL,
  `username` varchar(64) NOT NULL DEFAULT '',
  `attribute` varchar(64) NOT NULL DEFAULT '',
  `op` char(2) NOT NULL DEFAULT '==',
  `value` varchar(253) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `radcheck`
--

INSERT INTO `radcheck` (`id`, `username`, `attribute`, `op`, `value`) VALUES
(1, '4iml1u3r', 'Cleartext-Password', ':=', '4iml1u3r'),
(2, '4iml1u3r', 'User-Profile', ':=', '2M_Profile'),
(3, 'b31ltch4', 'Cleartext-Password', ':=', 'b31ltch4'),
(4, 'b31ltch4', 'User-Profile', ':=', '2M_Profile'),
(5, 'b38rn5mc', 'Cleartext-Password', ':=', 'b38rn5mc'),
(6, 'b38rn5mc', 'User-Profile', ':=', '2M_Profile'),
(7, '7wq13dyi', 'Cleartext-Password', ':=', '7wq13dyi'),
(8, '7wq13dyi', 'User-Profile', ':=', '2M_Profile');

-- --------------------------------------------------------

--
-- Struktur dari tabel `radgroupcheck`
--

CREATE TABLE `radgroupcheck` (
  `id` int UNSIGNED NOT NULL,
  `groupname` varchar(64) NOT NULL DEFAULT '',
  `attribute` varchar(64) NOT NULL DEFAULT '',
  `op` char(2) NOT NULL DEFAULT '==',
  `value` varchar(253) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `radgroupcheck`
--

INSERT INTO `radgroupcheck` (`id`, `groupname`, `attribute`, `op`, `value`) VALUES
(1, 'pppoe_1', 'Framed-Protocol', '==', 'PPP'),
(2, 'pppoe_2', 'Framed-Protocol', '==', 'PPP'),
(3, 'pppoe_3', 'Framed-Protocol', '==', 'PPP'),
(4, 'pppoe_4', 'Framed-Protocol', '==', 'PPP'),
(5, 'pppoe_4', 'Simultaneous-Use', ':=', '1'),
(6, 'pppoe_3', 'Simultaneous-Use', ':=', '1'),
(7, 'pppoe_1', 'Simultaneous-Use', ':=', '1'),
(8, 'pppoe_1', 'Simultaneous-Use', ':=', '1'),
(9, '2M_Profile', 'Simultaneous-Use', ':=', '1'),
(10, '3M_Profile', 'Simultaneous-Use', ':=', '1'),
(11, '5M_Profile', 'Simultaneous-Use', ':=', '1'),
(12, '7M_Profile', 'Simultaneous-Use', ':=', '1'),
(13, '10M_Profile', 'Simultaneous-Use', ':=', '1'),
(14, 'game', 'Simultaneous-Use', ':=', '1');

-- --------------------------------------------------------

--
-- Struktur dari tabel `radgroupreply`
--

CREATE TABLE `radgroupreply` (
  `id` int UNSIGNED NOT NULL,
  `groupname` varchar(64) NOT NULL DEFAULT '',
  `attribute` varchar(64) NOT NULL DEFAULT '',
  `op` char(2) NOT NULL DEFAULT '=',
  `value` varchar(253) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `radgroupreply`
--

INSERT INTO `radgroupreply` (`id`, `groupname`, `attribute`, `op`, `value`) VALUES
(1, '2M_Profile', 'Framed-Pool', '=', '2M_Pool'),
(2, '3M_Profile', 'Framed-Pool', '=', '3M_Pool'),
(3, '5M_Profile', 'Framed-Pool', '=', '5M_Pool'),
(4, '7M_Profile', 'Framed-Pool', '=', '7M_Pool'),
(5, '10M_Profile', 'Framed-Pool', '=', '10M_Pool'),
(6, 'game', 'Framed-Pool', '=', 'GAME_Pool'),
(7, 'pppoe_1', 'Framed-Pool', '=', 'PPPOE1_Pool'),
(8, 'pppoe_2', 'Framed-Pool', '=', 'PPPOE2_Pool'),
(9, 'pppoe_3', 'Framed-Pool', '=', 'PPPOE3_Pool'),
(10, 'pppoe_4', 'Framed-Pool', '=', 'PPPOE4_Pool'),
(11, 'pppoe_4', 'Mikrotik-Rate-Limit', '=', '2M/10M'),
(12, 'pppoe_3', 'Mikrotik-Rate-Limit', '=', '2M/7M'),
(13, 'pppoe_2', 'Mikrotik-Rate-Limit', '=', '2M/5M'),
(14, 'pppoe_1', 'Mikrotik-Rate-Limit', '=', '2M/3M'),
(15, '2M_Profile', 'Mikrotik-Rate-Limit', '=', '1M/2M'),
(16, '3M_Profile', 'Mikrotik-Rate-Limit', '=', '2M/3M'),
(17, '5M_Profile', 'Mikrotik-Rate-Limit', '=', '2M/5M'),
(18, '7M_Profile', 'Mikrotik-Rate-Limit', '=', '2M/7M'),
(19, '10M_Profile', 'Mikrotik-Rate-Limit', '=', '2M/10M'),
(20, 'game', 'Mikrotik-Rate-Limit', '=', '128k/128k');

-- --------------------------------------------------------

--
-- Struktur dari tabel `radpackage`
--

CREATE TABLE `radpackage` (
  `id` int NOT NULL,
  `id_package` int DEFAULT NULL,
  `username` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `radpackage`
--

INSERT INTO `radpackage` (`id`, `id_package`, `username`) VALUES
(1, 1, '4iml1u3r'),
(2, 1, 'b31ltch4'),
(3, 1, 'b38rn5mc'),
(4, 1, '7wq13dyi');

-- --------------------------------------------------------

--
-- Struktur dari tabel `radpostauth`
--

CREATE TABLE `radpostauth` (
  `id` int NOT NULL,
  `username` varchar(64) NOT NULL DEFAULT '',
  `pass` varchar(64) NOT NULL DEFAULT '',
  `reply` varchar(32) NOT NULL DEFAULT '',
  `authdate` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `radreply`
--

CREATE TABLE `radreply` (
  `id` int UNSIGNED NOT NULL,
  `username` varchar(64) NOT NULL DEFAULT '',
  `attribute` varchar(64) NOT NULL DEFAULT '',
  `op` char(2) NOT NULL DEFAULT '=',
  `value` varchar(253) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `radusergroup`
--

CREATE TABLE `radusergroup` (
  `username` varchar(64) NOT NULL DEFAULT '',
  `groupname` varchar(64) NOT NULL DEFAULT '',
  `priority` int NOT NULL DEFAULT '1'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `radusergroup`
--

INSERT INTO `radusergroup` (`username`, `groupname`, `priority`) VALUES
('2M_Profile', '2M_Profile', 8),
('5M_Profile', '5M_Profile', 8),
('3M_Profile', '3M_Profile', 8),
('7M_Profile', '7M_Profile', 8),
('10M_Profile', '10M_Profile', 1),
('pppoe_1', 'pppoe_1', 1),
('pppoe_2', 'pppoe_2', 1),
('pppoe_3', 'pppoe_3', 1),
('pppoe_4', 'pppoe_4', 1),
('game', 'game', 1);

-- --------------------------------------------------------

--
-- Struktur dari tabel `reseller`
--

CREATE TABLE `reseller` (
  `id` int NOT NULL,
  `telegram_id` int NOT NULL,
  `chat_id` int NOT NULL,
  `name` varchar(255) NOT NULL,
  `register_code` varchar(255) NOT NULL,
  `active` enum('yes','no') NOT NULL DEFAULT 'no',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `transaction`
--

CREATE TABLE `transaction` (
  `id` int NOT NULL,
  `id_reseller` int NOT NULL,
  `id_radpackage` int DEFAULT NULL,
  `transaction_code` varchar(255) NOT NULL,
  `status` enum('in','out') NOT NULL DEFAULT 'out',
  `value` int NOT NULL,
  `information` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='value';

--
-- Dumping data untuk tabel `transaction`
--

INSERT INTO `transaction` (`id`, `id_reseller`, `id_radpackage`, `transaction_code`, `status`, `value`, `information`, `created_at`) VALUES
(1, 1, NULL, '123104959456', 'in', 200000, 'okelah kalau begitu', '2021-04-26 07:56:07'),
(2, 1, 1, '1619425328906766', 'out', 2000, NULL, '2021-04-26 08:22:08'),
(3, 1, 2, '1619440036212883', 'out', 2000, NULL, '2021-04-26 12:27:16'),
(4, 1, 3, '1619440348485136', 'out', 2000, NULL, '2021-04-26 12:32:28'),
(5, 1, 4, '1619440780111313', 'out', 2000, NULL, '2021-04-26 12:39:40');

-- --------------------------------------------------------

--
-- Struktur dari tabel `users`
--

CREATE TABLE `users` (
  `id` int NOT NULL,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `level` enum('admin','user') NOT NULL DEFAULT 'user',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data untuk tabel `users`
--

INSERT INTO `users` (`id`, `username`, `password`, `level`, `created_at`) VALUES
(1, 'rhein', '$2a$04$LdiQ0FRd6HMuAi./7HZWIufQRlpW2YQd/dz6xDnckHJEnX1AY/ZEu', 'admin', '2021-04-25 16:21:38');

--
-- Indexes for dumped tables
--

--
-- Indeks untuk tabel `message`
--
ALTER TABLE `message`
  ADD PRIMARY KEY (`id`);

--
-- Indeks untuk tabel `nas`
--
ALTER TABLE `nas`
  ADD PRIMARY KEY (`id`),
  ADD KEY `nasname` (`nasname`);

--
-- Indeks untuk tabel `package`
--
ALTER TABLE `package`
  ADD PRIMARY KEY (`id`),
  ADD KEY `profile` (`profile`);

--
-- Indeks untuk tabel `radacct`
--
ALTER TABLE `radacct`
  ADD PRIMARY KEY (`radacctid`),
  ADD UNIQUE KEY `acctuniqueid` (`acctuniqueid`),
  ADD KEY `username` (`username`),
  ADD KEY `framedipaddress` (`framedipaddress`),
  ADD KEY `acctsessionid` (`acctsessionid`),
  ADD KEY `acctsessiontime` (`acctsessiontime`),
  ADD KEY `acctstarttime` (`acctstarttime`),
  ADD KEY `acctinterval` (`acctinterval`),
  ADD KEY `acctstoptime` (`acctstoptime`),
  ADD KEY `nasipaddress` (`nasipaddress`);

--
-- Indeks untuk tabel `radcheck`
--
ALTER TABLE `radcheck`
  ADD PRIMARY KEY (`id`),
  ADD KEY `username` (`username`(32));

--
-- Indeks untuk tabel `radgroupcheck`
--
ALTER TABLE `radgroupcheck`
  ADD PRIMARY KEY (`id`),
  ADD KEY `groupname` (`groupname`(32));

--
-- Indeks untuk tabel `radgroupreply`
--
ALTER TABLE `radgroupreply`
  ADD PRIMARY KEY (`id`),
  ADD KEY `groupname` (`groupname`(32));

--
-- Indeks untuk tabel `radpackage`
--
ALTER TABLE `radpackage`
  ADD PRIMARY KEY (`id`),
  ADD KEY `id_package` (`id_package`);

--
-- Indeks untuk tabel `radpostauth`
--
ALTER TABLE `radpostauth`
  ADD PRIMARY KEY (`id`);

--
-- Indeks untuk tabel `radreply`
--
ALTER TABLE `radreply`
  ADD PRIMARY KEY (`id`),
  ADD KEY `username` (`username`(32));

--
-- Indeks untuk tabel `radusergroup`
--
ALTER TABLE `radusergroup`
  ADD KEY `username` (`username`(32));

--
-- Indeks untuk tabel `transaction`
--
ALTER TABLE `transaction`
  ADD PRIMARY KEY (`id`),
  ADD KEY `id_reseller` (`id_reseller`),
  ADD KEY `id_radpackage` (`id_radpackage`);

--
-- Indeks untuk tabel `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT untuk tabel yang dibuang
--

--
-- AUTO_INCREMENT untuk tabel `message`
--
ALTER TABLE `message`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `nas`
--
ALTER TABLE `nas`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `package`
--
ALTER TABLE `package`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT untuk tabel `radacct`
--
ALTER TABLE `radacct`
  MODIFY `radacctid` bigint NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `radcheck`
--
ALTER TABLE `radcheck`
  MODIFY `id` int UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT untuk tabel `radgroupcheck`
--
ALTER TABLE `radgroupcheck`
  MODIFY `id` int UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=15;

--
-- AUTO_INCREMENT untuk tabel `radgroupreply`
--
ALTER TABLE `radgroupreply`
  MODIFY `id` int UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=21;

--
-- AUTO_INCREMENT untuk tabel `radpackage`
--
ALTER TABLE `radpackage`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT untuk tabel `radpostauth`
--
ALTER TABLE `radpostauth`
  MODIFY `id` int NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `radreply`
--
ALTER TABLE `radreply`
  MODIFY `id` int UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `transaction`
--
ALTER TABLE `transaction`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT untuk tabel `users`
--
ALTER TABLE `users`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
