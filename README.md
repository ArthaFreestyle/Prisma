# 🏆 Prisma (Prestasi Mahasiswa)
Sistem Pelaporan Prestasi Mahasiswa yang memudahkan pengelolaan data prestasi, baik di PostgreSQL maupun MongoDB.

---

## 📂 Database Migration

Semua migration database disimpan di folder `db/migrations_*` sesuai tipe database:  

- `db/migrations_postgre` → PostgreSQL  
- `db/migrations_mongo` → MongoDB  

---
## ➕ Create Migration
postgre
```bash
    migrate create -ext sql -dir db/migrations_postgre create_table_xxx
```

mongodb
```bash
    migrate create -ext js -dir db/migrations_mongo create_collection_xxx
```

## ▶️ Run Migration

postgre
```bash
./migratepg.sh
```