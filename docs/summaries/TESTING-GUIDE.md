# 🧪 Testing Guide - CAKLI v1

This guide provides step-by-step instructions for testing all features of the CAKLI admin dashboard.

---

## 🚀 Prerequisites

Before testing, ensure all services are running:

```bash
# Start backend services (PostgreSQL, MinIO, API)
make dev

# In another terminal, start frontend
cd apps/web
npm run dev
```

**Access URLs:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- MinIO Console: http://localhost:9001

**Default Admin Credentials:**
- Email: `master@cakli.id`
- Password: `admin123`

---

## 1. Authentication Testing

### 1.1 Login Flow

**Test Case:** Successful login
1. Navigate to http://localhost:3000
2. Enter email: `master@cakli.id`
3. Enter password: `admin123`
4. Click "Masuk"
5. **Expected:** Redirect to `/operation-admin` dashboard
6. **Expected:** See admin name in sidebar

**Test Case:** Invalid credentials
1. Enter email: `wrong@email.com`
2. Enter password: `wrongpassword`
3. Click "Masuk"
4. **Expected:** Error toast "Email atau password salah"

**Test Case:** Empty fields
1. Leave email and password empty
2. Click "Masuk"
3. **Expected:** Browser validation errors

### 1.2 Protected Routes

**Test Case:** Access protected route without login
1. Open new incognito window
2. Navigate to http://localhost:3000/operation-admin
3. **Expected:** Redirect to login page

**Test Case:** Token refresh
1. Login successfully
2. Wait 15+ minutes (or modify JWT expiry for faster testing)
3. Perform any action (e.g., navigate to drivers page)
4. **Expected:** Token auto-refreshes, no logout

### 1.3 Logout Flow

**Test Case:** Successful logout
1. Login successfully
2. Click on admin avatar in sidebar
3. Click "Logout"
4. **Expected:** Success toast "Berhasil logout"
5. **Expected:** Redirect to login page
6. **Expected:** Cannot access protected routes

---

## 2. Dashboard Testing

### 2.1 Dashboard Statistics

**Test Case:** View dashboard stats
1. Login and navigate to `/operation-admin`
2. **Expected:** See stat cards:
   - Total Driver (count from database)
   - Drivers Online (count from database)
   - Active Complaints (count from database)
3. **Expected:** Stats match database counts

**Test Case:** Real-time updates
1. Open dashboard
2. In another tab, create a new driver via API or form
3. Refresh dashboard
4. **Expected:** Driver count increases

---

## 3. Driver Management Testing

### 3.1 Driver List

**Test Case:** View driver list
1. Navigate to `/operation-admin/drivers`
2. **Expected:** See table with columns:
   - Driver (photo + name + email)
   - NIK
   - Kontak (phone)
   - Bank (name + account number)
   - Status (badge)
   - Aksi (actions)
3. **Expected:** Pagination controls at bottom

**Test Case:** Search drivers
1. Enter "test" in search box
2. Wait for debounce (500ms)
3. **Expected:** Table filters to show only matching drivers
4. **Expected:** Pagination resets to page 1

**Test Case:** Filter by status
1. Select "Verified" from status dropdown
2. **Expected:** Table shows only drivers with status "accepted"
3. Select "Pending"
4. **Expected:** Table shows only drivers with status "pending"
5. Select "Semua Status"
6. **Expected:** Table shows all drivers

**Test Case:** Pagination
1. If total drivers > 10, check pagination
2. Click "Next" button
3. **Expected:** Navigate to page 2
4. **Expected:** URL or state updates
5. Click "Previous" button
6. **Expected:** Navigate back to page 1

### 3.2 Create Driver

**Test Case:** Create driver with all fields
1. Click "Tambah Driver Baru" button
2. Fill in all fields:
   - Upload Foto Profile (JPG/PNG)
   - Upload Foto KTP (JPG/PNG)
   - Upload Foto Muka (JPG/PNG)
   - Nama Lengkap: "Test Driver"
   - Email: "testdriver@example.com"
   - Password: "password123"
   - No. Telepon: "081234567890"
   - NIK: "1234567890123456" (16 digits)
   - Tempat Lahir: "Jakarta"
   - Tanggal Lahir: "1990-01-01"
   - Bank: Select any bank
   - No. Rekening: "1234567890"
   - Status Verifikasi: "Pending"
3. Click "Simpan"
4. **Expected:** Loading state shows "Uploading..." then "Menyimpan..."
5. **Expected:** Success toast "Driver berhasil ditambahkan"
6. **Expected:** Modal closes
7. **Expected:** Table refreshes with new driver

**Test Case:** Create driver with missing required fields
1. Click "Tambah Driver Baru"
2. Leave required fields empty
3. Click "Simpan"
4. **Expected:** Error toast "Mohon lengkapi semua field yang wajib diisi"

**Test Case:** Create driver with invalid NIK
1. Click "Tambah Driver Baru"
2. Fill all fields but enter NIK with less than 16 digits
3. Click "Simpan"
4. **Expected:** Error toast "NIK harus 16 digit"

**Test Case:** Create driver with duplicate email
1. Click "Tambah Driver Baru"
2. Use email that already exists
3. Click "Simpan"
4. **Expected:** Error toast from backend about duplicate email

### 3.3 Edit Driver

**Test Case:** Edit driver
1. Click "..." menu on any driver row
2. Click "Edit Driver"
3. **Expected:** Modal opens with pre-filled data
4. Change name to "Updated Driver Name"
5. Leave password empty (should not update password)
6. Click "Simpan"
7. **Expected:** Success toast "Driver berhasil diupdate"
8. **Expected:** Table shows updated name

**Test Case:** Edit driver with new photo
1. Click "..." menu on any driver
2. Click "Edit Driver"
3. Upload new Foto Profile
4. **Expected:** Preview shows new image
5. Click "Simpan"
6. **Expected:** Photo uploads successfully
7. **Expected:** Driver updated with new photo

### 3.4 File Upload

**Test Case:** Upload valid image
1. Click "Tambah Driver Baru"
2. Click on "Foto Profile" upload area
3. Select a JPG or PNG file (< 5MB)
4. **Expected:** Preview shows uploaded image
5. **Expected:** Remove button (X) appears

**Test Case:** Remove uploaded image
1. Upload an image
2. Click X button on preview
3. **Expected:** Preview clears
4. **Expected:** Upload area shows again

**Test Case:** Upload invalid file type
1. Try to select a PDF or other non-image file
2. **Expected:** File picker only shows JPG/PNG files

---

## 4. User Management Testing

### 4.1 User List

**Test Case:** View user list
1. Navigate to `/operation-admin/users`
2. **Expected:** See table with columns:
   - User (avatar + name + ID)
   - Email
   - Telepon
   - Status (Active/Inactive badge)
   - Bergabung (join date)
   - Aksi (actions)

**Test Case:** Search users
1. Enter search query
2. **Expected:** Table filters to matching users
3. **Expected:** Pagination resets

**Test Case:** Pagination
1. Navigate through pages
2. **Expected:** Pagination works correctly

### 4.2 Create User

**Test Case:** Create user with all fields
1. Click "Tambah User Baru"
2. Fill in:
   - Upload Foto Profile
   - Nama Lengkap: "Test User"
   - Email: "testuser@example.com"
   - No. Telepon: "081234567890"
   - Password: "password123"
3. Click "Simpan"
4. **Expected:** Success toast "User berhasil ditambahkan"
5. **Expected:** Table refreshes with new user

**Test Case:** Create user without password
1. Click "Tambah User Baru"
2. Fill all fields except password
3. Click "Simpan"
4. **Expected:** Error toast "Password wajib diisi untuk user baru"

### 4.3 Edit User

**Test Case:** Edit user without changing password
1. Click "..." menu on any user
2. Click "Edit User"
3. Change name
4. Leave password empty
5. Click "Simpan"
6. **Expected:** Success toast "User berhasil diupdate"
7. **Expected:** Password not changed (can still login with old password)

**Test Case:** Edit user with new password
1. Edit user
2. Enter new password
3. Click "Simpan"
4. **Expected:** Password updated successfully

---

## 5. API Endpoint Testing

### 5.1 Using cURL

**Test Authentication:**
```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"master@cakli.id","password":"admin123"}' \
  -c cookies.txt

# Get current admin
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Logout
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -b cookies.txt
```

**Test Driver Endpoints:**
```bash
# List drivers
curl -X GET "http://localhost:8080/api/v1/admin/drivers?page=1&limit=10" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Get driver detail
curl -X GET http://localhost:8080/api/v1/admin/drivers/DRIVER_ID \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"

# Create driver
curl -X POST http://localhost:8080/api/v1/admin/drivers \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Driver",
    "email": "test@example.com",
    "password": "password123",
    "phone": "081234567890",
    "nik": "1234567890123456"
  }'
```

**Test Upload Endpoints:**
```bash
# Generate presigned upload URL
curl -X POST http://localhost:8080/api/v1/upload/presigned-url \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "filename": "test.jpg",
    "content_type": "image/jpeg",
    "folder": "drivers/photo-profile"
  }'

# Upload file to presigned URL (use URL from previous response)
curl -X PUT "PRESIGNED_UPLOAD_URL" \
  -H "Content-Type: image/jpeg" \
  --data-binary @test.jpg

# Generate presigned view URL
curl -X GET "http://localhost:8080/api/v1/upload/presigned-view-url?key=OBJECT_KEY" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 5.2 Using Postman

1. Import the API collection (if available)
2. Set environment variables:
   - `base_url`: http://localhost:8080/api/v1
   - `access_token`: (will be set after login)
3. Test each endpoint in order:
   - Auth → Login
   - Auth → Me
   - Dashboard → Stats
   - Drivers → List
   - Drivers → Create
   - Drivers → Update
   - Users → List
   - Users → Create
   - Banks → List
   - Upload → Presigned URL

---

## 6. Error Handling Testing

### 6.1 Network Errors

**Test Case:** Backend offline
1. Stop backend API (`docker stop cakli-api`)
2. Try to login
3. **Expected:** Error toast "Network error" or similar
4. Start backend again

**Test Case:** Slow network
1. Use browser DevTools to throttle network
2. Perform actions
3. **Expected:** Loading states show appropriately
4. **Expected:** Actions complete when network recovers

### 6.2 Validation Errors

**Test Case:** Invalid email format
1. Try to create driver with email "notanemail"
2. **Expected:** Validation error

**Test Case:** Duplicate email
1. Create driver with existing email
2. **Expected:** Backend error toast

**Test Case:** Invalid file type
1. Try to upload non-image file
2. **Expected:** Validation error

### 6.3 Authorization Errors

**Test Case:** Expired token
1. Manually expire access token (or wait 15 min)
2. Perform action
3. **Expected:** Token auto-refreshes
4. **Expected:** Action completes successfully

**Test Case:** Invalid token
1. Manually set invalid token in localStorage
2. Perform action
3. **Expected:** Redirect to login

---

## 7. Performance Testing

### 7.1 Load Testing

**Test Case:** Large dataset
1. Insert 1000+ drivers into database
2. Navigate to drivers page
3. **Expected:** Page loads within 2 seconds
4. **Expected:** Pagination works smoothly

**Test Case:** Search performance
1. With large dataset, perform search
2. **Expected:** Results appear within 1 second
3. **Expected:** Debounce prevents excessive API calls

### 7.2 File Upload Performance

**Test Case:** Large file upload
1. Try to upload 5MB image
2. **Expected:** Upload completes successfully
3. **Expected:** Progress indicator shows (if implemented)

**Test Case:** Multiple file uploads
1. Create driver with 3 photos
2. **Expected:** All photos upload sequentially
3. **Expected:** Loading state shows "Uploading..."

---

## 8. Browser Compatibility Testing

Test the application in:
- ✅ Chrome (latest)
- ✅ Firefox (latest)
- ✅ Safari (latest)
- ✅ Edge (latest)

**Test Case:** Cross-browser functionality
1. Test all features in each browser
2. **Expected:** All features work consistently
3. **Expected:** UI renders correctly

---

## 9. Responsive Design Testing

### 9.1 Mobile View (< 768px)

**Test Case:** Mobile navigation
1. Resize browser to mobile width
2. **Expected:** Sidebar collapses or becomes hamburger menu
3. **Expected:** Tables scroll horizontally if needed

**Test Case:** Mobile forms
1. Open create driver modal on mobile
2. **Expected:** Form is usable and scrollable
3. **Expected:** Buttons are accessible

### 9.2 Tablet View (768px - 1024px)

**Test Case:** Tablet layout
1. Resize to tablet width
2. **Expected:** Layout adapts appropriately
3. **Expected:** All features accessible

---

## 10. Security Testing

### 10.1 Authentication Security

**Test Case:** Password hashing
1. Check database
2. **Expected:** Passwords are bcrypt hashed, not plain text

**Test Case:** Token security
1. Check localStorage
2. **Expected:** Only access token stored
3. **Expected:** Refresh token in httpOnly cookie

**Test Case:** CORS protection
1. Try to access API from different origin
2. **Expected:** CORS headers allow only configured origins

### 10.2 Authorization Security

**Test Case:** Protected endpoints
1. Try to access `/api/v1/admin/drivers` without token
2. **Expected:** 401 Unauthorized

**Test Case:** SQL injection
1. Try to inject SQL in search field: `'; DROP TABLE drivers; --`
2. **Expected:** Query safely escaped, no SQL injection

---

## 11. Regression Testing Checklist

After any code changes, verify:

- [ ] Login/logout works
- [ ] Dashboard stats load correctly
- [ ] Driver list loads with pagination
- [ ] Driver search works
- [ ] Driver create works with file upload
- [ ] Driver edit works
- [ ] User list loads with pagination
- [ ] User search works
- [ ] User create works with file upload
- [ ] User edit works
- [ ] Token refresh works
- [ ] Error handling works
- [ ] No console errors
- [ ] No broken images
- [ ] Responsive design works

---

## 12. Known Issues & Workarounds

### Issue 1: Docker port conflict
**Symptom:** `Bind for 0.0.0.0:9000 failed: port is already allocated`
**Solution:** See `docs/TROUBLESHOOTING.md`

### Issue 2: Database connection fails
**Symptom:** `connection refused`
**Solution:** Ensure PostgreSQL container is running: `docker ps`

### Issue 3: MinIO bucket not found
**Symptom:** `bucket does not exist`
**Solution:** Restart API container, it will auto-create bucket

---

## 13. Test Data

### Sample Driver Data
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "password": "password123",
  "phone": "081234567890",
  "nik": "3201234567890123",
  "birth_place": "Jakarta",
  "birth_date": "1990-01-15",
  "bank_id": 1,
  "bank_account_number": "1234567890",
  "verification_status": "pending"
}
```

### Sample User Data
```json
{
  "name": "Jane Smith",
  "email": "jane.smith@example.com",
  "password": "password123",
  "phone": "081234567891"
}
```

---

## 14. Automated Testing (Future)

### Backend Tests (Go)
```bash
cd apps/api
go test ./... -v
```

### Frontend Tests (Jest + React Testing Library)
```bash
cd apps/web
npm test
```

### E2E Tests (Playwright/Cypress)
```bash
cd apps/web
npm run test:e2e
```

---

## 📝 Test Report Template

After testing, document results:

```markdown
# Test Report - CAKLI v1

**Date:** YYYY-MM-DD
**Tester:** [Name]
**Environment:** Development

## Test Results

| Test Case | Status | Notes |
|-----------|--------|-------|
| Login | ✅ Pass | - |
| Logout | ✅ Pass | - |
| Dashboard Stats | ✅ Pass | - |
| Driver List | ✅ Pass | - |
| Driver Create | ❌ Fail | File upload timeout |
| Driver Edit | ✅ Pass | - |
| User List | ✅ Pass | - |
| User Create | ✅ Pass | - |
| User Edit | ✅ Pass | - |

## Issues Found

1. **File upload timeout** - Large files (>5MB) fail to upload
   - Severity: Medium
   - Steps to reproduce: Upload 10MB image
   - Expected: Success
   - Actual: Timeout after 30s

## Recommendations

1. Increase upload timeout
2. Add file size validation
3. Add progress indicator
```

---

**Happy Testing! 🎉**
