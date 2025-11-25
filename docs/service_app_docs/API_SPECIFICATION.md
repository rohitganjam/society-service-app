# API Specification Document

**Version:** 1.0
**Date:** November 20, 2025
**Purpose:** Complete API specifications for all service flows

---

## Table of Contents

1. [Authentication APIs](#1-authentication-apis)
2. [Onboarding APIs](#2-onboarding-apis)
   - [Resident Onboarding](#21-resident-onboarding)
   - [Vendor Onboarding](#22-vendor-onboarding)
3. [Approval & Verification APIs](#3-approval--verification-apis)
   - [Society Admin Approvals](#31-society-admin-approvals)
     - [Resident Approvals](#311-get-pending-resident-approvals)
     - [Vendor Approvals](#313-get-pending-vendor-approvals)
     - [Group Management (Unified Buildings/Phases)](#315-manage-society-groups-unified-for-buildingsphases)
     - [Vendor Service Area Assignment](#316-assign-vendor-to-service-areas-unified-groups)
   - [Platform Admin Approvals](#32-platform-admin-approvals)
4. [Rate Card Management APIs](#4-rate-card-management-apis)
5. [Vendor Listing & Discovery APIs](#5-vendor-listing--discovery-apis)
6. [Order Management APIs](#6-order-management-apis)
7. [Workflow Management APIs](#7-workflow-management-apis)
8. [Payment APIs](#8-payment-apis)

---

## 1. Authentication APIs

### 1.1 Send OTP (Phone)

**Endpoint:** `POST /api/v1/auth/send-otp`

**Description:** Send OTP to user's phone number for verification (login or registration)

**Request Body:**
```json
{
  "phone": "+919876543210"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "OTP sent successfully",
  "data": {
    "otp_id": "uuid-v4",
    "phone": "+919876543210",
    "expires_at": "2025-11-20T10:15:00Z",
    "retry_after": 60
  }
}
```

**Response (429 Too Many Requests):**
```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Too many OTP requests. Please try again after 60 seconds",
    "retry_after": 60
  }
}
```

---

### 1.2 Send OTP (Email)

**Endpoint:** `POST /api/v1/auth/send-email-otp`

**Description:** Send OTP to user's email for verification (login or registration)

**Request Body:**
```json
{
  "email": "ramesh@example.com"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "OTP sent to email successfully",
  "data": {
    "otp_id": "uuid-v4",
    "email": "r****h@example.com",
    "expires_at": "2025-11-20T10:15:00Z",
    "retry_after": 60
  }
}
```

---

### 1.3 Verify OTP (Phone)

**Endpoint:** `POST /api/v1/auth/verify-otp`

**Description:** Verify phone OTP and authenticate user

**Request Body:**
```json
{
  "phone": "+919876543210",
  "otp": "123456",
  "otp_id": "uuid-v4"
}
```

**Response (200 OK - Existing User):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "user_id": "uuid-v4",
      "phone": "+919876543210",
      "email": "ramesh@example.com",
      "full_name": "Ramesh Kumar",
      "user_type": "RESIDENT",
      "is_verified": true,
      "profile_photo_url": "https://...",
      "has_password": false
    },
    "tokens": {
      "access_token": "jwt-token",
      "refresh_token": "jwt-refresh-token",
      "expires_in": 3600
    },
    "is_new_user": false
  }
}
```

**Response (200 OK - New User):**
```json
{
  "success": true,
  "message": "OTP verified. Please complete registration",
  "data": {
    "temp_user_id": "uuid-v4",
    "phone": "+919876543210",
    "is_new_user": true,
    "requires_registration": true
  }
}
```

**Response (401 Unauthorized):**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_OTP",
    "message": "Invalid or expired OTP"
  }
}
```

---

### 1.4 Verify OTP (Email)

**Endpoint:** `POST /api/v1/auth/verify-email-otp`

**Description:** Verify email OTP and authenticate user

**Request Body:**
```json
{
  "email": "ramesh@example.com",
  "otp": "123456",
  "otp_id": "uuid-v4"
}
```

**Response (200 OK - Existing User):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "user_id": "uuid-v4",
      "phone": "+919876543210",
      "email": "ramesh@example.com",
      "full_name": "Ramesh Kumar",
      "user_type": "RESIDENT",
      "is_verified": true,
      "profile_photo_url": "https://...",
      "has_password": false
    },
    "tokens": {
      "access_token": "jwt-token",
      "refresh_token": "jwt-refresh-token",
      "expires_in": 3600
    },
    "is_new_user": false
  }
}
```

**Response (200 OK - New User):**
```json
{
  "success": true,
  "message": "OTP verified. Please complete registration",
  "data": {
    "temp_user_id": "uuid-v4",
    "email": "ramesh@example.com",
    "is_new_user": true,
    "requires_registration": true
  }
}
```

---

### 1.5 Login with Password

**Endpoint:** `POST /api/v1/auth/login`

**Description:** Login using email/phone and password (if user has set password)

**Request Body (Email):**
```json
{
  "email": "ramesh@example.com",
  "password": "SecurePassword123!"
}
```

**Request Body (Phone):**
```json
{
  "phone": "+919876543210",
  "password": "SecurePassword123!"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "user_id": "uuid-v4",
      "phone": "+919876543210",
      "email": "ramesh@example.com",
      "full_name": "Ramesh Kumar",
      "user_type": "RESIDENT",
      "is_verified": true,
      "has_password": true
    },
    "tokens": {
      "access_token": "jwt-token",
      "refresh_token": "jwt-refresh-token",
      "expires_in": 3600
    }
  }
}
```

**Response (401 Unauthorized):**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid email/phone or password"
  }
}
```

---

### 1.6 OAuth Login (Google)

**Endpoint:** `POST /api/v1/auth/oauth/google`

**Description:** Authenticate using Google OAuth

**Request Body:**
```json
{
  "id_token": "google-id-token",
  "access_token": "google-access-token"
}
```

**Response (200 OK - Existing User):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "user_id": "uuid-v4",
      "email": "ramesh@example.com",
      "full_name": "Ramesh Kumar",
      "user_type": "RESIDENT",
      "is_verified": true,
      "profile_photo_url": "https://...",
      "oauth_provider": "GOOGLE"
    },
    "tokens": {
      "access_token": "jwt-token",
      "refresh_token": "jwt-refresh-token",
      "expires_in": 3600
    },
    "is_new_user": false
  }
}
```

**Response (200 OK - New User):**
```json
{
  "success": true,
  "message": "OAuth verified. Please complete registration",
  "data": {
    "temp_user_id": "uuid-v4",
    "email": "ramesh@example.com",
    "full_name": "Ramesh Kumar",
    "profile_photo_url": "https://...",
    "is_new_user": true,
    "requires_registration": true,
    "oauth_provider": "GOOGLE"
  }
}
```

---

### 1.7 OAuth Login (Facebook)

**Endpoint:** `POST /api/v1/auth/oauth/facebook`

**Description:** Authenticate using Facebook OAuth

**Request Body:**
```json
{
  "access_token": "facebook-access-token"
}
```

**Response:** Same structure as Google OAuth

---

### 1.8 Forgot Password

**Endpoint:** `POST /api/v1/auth/forgot-password`

**Description:** Request password reset OTP

**Request Body (Email):**
```json
{
  "email": "ramesh@example.com"
}
```

**Request Body (Phone):**
```json
{
  "phone": "+919876543210"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password reset OTP sent successfully",
  "data": {
    "otp_id": "uuid-v4",
    "sent_to": "r****h@example.com",
    "expires_at": "2025-11-20T10:15:00Z"
  }
}
```

**Response (404 Not Found):**
```json
{
  "success": false,
  "error": {
    "code": "USER_NOT_FOUND",
    "message": "No account found with this email/phone"
  }
}
```

---

### 1.9 Reset Password

**Endpoint:** `POST /api/v1/auth/reset-password`

**Description:** Reset password using OTP

**Request Body:**
```json
{
  "email": "ramesh@example.com",
  "otp": "123456",
  "otp_id": "uuid-v4",
  "new_password": "NewSecurePassword123!"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password reset successfully",
  "data": {
    "user_id": "uuid-v4",
    "message": "You can now login with your new password"
  }
}
```

**Response (401 Unauthorized):**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_OTP",
    "message": "Invalid or expired OTP"
  }
}
```

---

### 1.10 Set Password

**Endpoint:** `POST /api/v1/auth/set-password`

**Description:** Set password for account that doesn't have one (e.g., OAuth or OTP-only users)

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body:**
```json
{
  "password": "SecurePassword123!",
  "confirm_password": "SecurePassword123!"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password set successfully",
  "data": {
    "user_id": "uuid-v4",
    "has_password": true
  }
}
```

---

### 1.11 Change Password

**Endpoint:** `POST /api/v1/auth/change-password`

**Description:** Change existing password (requires current password)

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body:**
```json
{
  "current_password": "OldPassword123!",
  "new_password": "NewPassword123!",
  "confirm_password": "NewPassword123!"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password changed successfully"
}
```

**Response (401 Unauthorized):**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_PASSWORD",
    "message": "Current password is incorrect"
  }
}
```

---

### 1.12 Refresh Token

**Endpoint:** `POST /api/v1/auth/refresh`

**Description:** Refresh access token using refresh token

**Request Body:**
```json
{
  "refresh_token": "jwt-refresh-token"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "access_token": "new-jwt-token",
    "expires_in": 3600
  }
}
```

---

### 1.13 Logout

**Endpoint:** `POST /api/v1/auth/logout`

**Description:** Logout user and invalidate refresh token

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body:**
```json
{
  "refresh_token": "jwt-refresh-token"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

---

## 2. Onboarding APIs

### 2.1 Resident Onboarding

#### 2.1.1 Check Society Roster

**Endpoint:** `POST /api/v1/onboarding/resident/check-roster`

**Description:** Check if phone number exists in society roster(s). Supports multiple society memberships and independent houses with multiple households.

**Request Body:**
```json
{
  "phone": "+919876543210"
}
```

**Response (200 OK - Found in Single Society - Apartment):**
```json
{
  "success": true,
  "data": {
    "found_in_roster": true,
    "has_multiple_societies": false,
    "residences": [
      {
        "society_id": 1,
        "society_name": "Maple Gardens",
        "society_type": "APARTMENT",
        "address": "123 MG Road, Koramangala",
        "city": "Bangalore",
        "flat_number": "A-404",
        "tower": "A",
        "floor": 4,
        "unit_type": "FLAT",
        "suggested_name": "Ramesh Kumar"
      }
    ]
  }
}
```

**Response (200 OK - Found in Single Society - Independent House):**
```json
{
  "success": true,
  "data": {
    "found_in_roster": true,
    "has_multiple_societies": false,
    "residences": [
      {
        "society_id": 5,
        "society_name": "Green Meadows Layout",
        "society_type": "LAYOUT",
        "address": "House #42, 5th Cross, Green Meadows",
        "city": "Bangalore",
        "house_number": "42",
        "street": "5th Cross",
        "floor": 1,
        "unit_type": "HOUSE",
        "suggested_name": "Priya Sharma",
        "notes": "Ground floor"
      }
    ]
  }
}
```

**Response (200 OK - Found in Multiple Societies):**
```json
{
  "success": true,
  "data": {
    "found_in_roster": true,
    "has_multiple_societies": true,
    "residences": [
      {
        "society_id": 1,
        "society_name": "Maple Gardens",
        "society_type": "APARTMENT",
        "address": "123 MG Road, Koramangala, Bangalore",
        "city": "Bangalore",
        "flat_number": "A-404",
        "tower": "A",
        "floor": 4,
        "unit_type": "FLAT",
        "suggested_name": "Ramesh Kumar",
        "is_primary": true
      },
      {
        "society_id": 3,
        "society_name": "Beach View Apartments",
        "society_type": "APARTMENT",
        "address": "456 Beach Road, Chennai",
        "city": "Chennai",
        "flat_number": "201",
        "floor": 2,
        "unit_type": "FLAT",
        "suggested_name": "Ramesh Kumar",
        "is_primary": false,
        "notes": "Weekend home"
      },
      {
        "society_id": 7,
        "society_name": "Hill View Layout",
        "society_type": "LAYOUT",
        "address": "House #15, Hill View Road, Ooty",
        "city": "Ooty",
        "house_number": "15",
        "floor": 2,
        "unit_type": "HOUSE",
        "suggested_name": "Ramesh Kumar",
        "is_primary": false,
        "notes": "First floor - Rented"
      }
    ],
    "message": "Multiple residences found. Please select which society to use as primary."
  }
}
```

**Response (200 OK - Not Found in Roster):**
```json
{
  "success": true,
  "data": {
    "found_in_roster": false,
    "message": "Phone number not found in any society roster. Please search and select your society."
  }
}
```

**Notes:**
- `society_type`: `APARTMENT` (multi-unit building) or `LAYOUT` (independent houses)
- `unit_type`: `FLAT` (apartment) or `HOUSE` (independent house)
- For independent houses with multiple households: Same `house_number` but different `floor` values
- For multi-society users: First residence where `is_primary: true` is used as default
- Independent houses may have `floor: 0` (ground), `floor: 1` (first floor), `floor: 2` (second floor), etc.

---

#### 2.1.2 Search Societies

**Endpoint:** `GET /api/v1/societies/search`

**Description:** Search for societies by name with autocomplete support. Used when resident is not found in roster.

**Query Parameters:**
- `q` (required): Search query (minimum 2 characters)
- `city` (optional): Filter by city
- `society_type` (optional): Filter by type (APARTMENT, LAYOUT)
- `limit` (optional): Number of results (default: 10, max: 50)

**Example Request:**
```
GET /api/v1/societies/search?q=maple&city=Bangalore
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "societies": [
      {
        "society_id": 1,
        "name": "Maple Gardens",
        "society_type": "APARTMENT",
        "address": "123 MG Road, Koramangala",
        "city": "Bangalore",
        "state": "Karnataka",
        "pincode": "560034",
        "total_flats": 250,
        "total_houses": null,
        "occupied_flats": 230,
        "status": "ACTIVE",
        "has_subscription": true
      },
      {
        "society_id": 8,
        "name": "Maple Grove Layout",
        "society_type": "LAYOUT",
        "address": "Maple Grove Road, Sarjapur",
        "city": "Bangalore",
        "state": "Karnataka",
        "pincode": "560035",
        "total_flats": null,
        "total_houses": 85,
        "occupied_houses": 78,
        "status": "ACTIVE",
        "has_subscription": true
      }
    ],
    "total_results": 2,
    "query": "maple"
  }
}
```

**Response (200 OK - No Results):**
```json
{
  "success": true,
  "data": {
    "societies": [],
    "total_results": 0,
    "query": "xyz",
    "message": "No societies found matching 'xyz'. Please check the spelling or contact support."
  }
}
```

**Notes:**
- Search is performed on society name, address, and pincode
- Results are sorted by relevance (name match first, then address match)
- Only active societies with valid subscriptions are returned
- For APARTMENT type: `total_flats` and `occupied_flats` are populated
- For LAYOUT type: `total_houses` and `occupied_houses` are populated

---

#### 2.1.3 Complete Resident Registration

**Endpoint:** `POST /api/v1/onboarding/resident/register`

**Description:** Complete resident registration. Supports apartments and independent houses with multiple households.

**Request Body (Auto-verified from roster - Apartment):**
```json
{
  "temp_user_id": "uuid-v4",
  "phone": "+919876543210",
  "full_name": "Ramesh Kumar",
  "email": "ramesh@example.com",
  "society_id": 1,
  "unit_type": "FLAT",
  "flat_number": "A-404",
  "tower": "A",
  "floor": 4,
  "from_roster": true
}
```

**Request Body (Auto-verified from roster - Independent House):**
```json
{
  "temp_user_id": "uuid-v4",
  "phone": "+919876543211",
  "full_name": "Priya Sharma",
  "email": "priya@example.com",
  "society_id": 5,
  "unit_type": "HOUSE",
  "house_number": "42",
  "street": "5th Cross",
  "floor": 1,
  "notes": "First floor",
  "from_roster": true
}
```

**Request Body (Manual verification - Not in roster - Apartment):**
```json
{
  "temp_user_id": "uuid-v4",
  "phone": "+919876543210",
  "full_name": "Ramesh Kumar",
  "email": "ramesh@example.com",
  "society_id": 1,
  "unit_type": "FLAT",
  "flat_number": "B-302",
  "tower": "B",
  "floor": 3,
  "from_roster": false
}
```

**Request Body (Manual verification - Not in roster - Independent House):**
```json
{
  "temp_user_id": "uuid-v4",
  "phone": "+919876543212",
  "full_name": "Amit Verma",
  "email": "amit@example.com",
  "society_id": 8,
  "unit_type": "HOUSE",
  "house_number": "15",
  "street": "Main Road",
  "floor": 0,
  "notes": "Ground floor - separate entrance",
  "from_roster": false
}
```

**Request Body (Multi-society user - Adding additional society):**
```json
{
  "user_id": "existing-uuid-v4",
  "phone": "+919876543210",
  "society_id": 3,
  "unit_type": "FLAT",
  "flat_number": "201",
  "floor": 2,
  "is_primary": false,
  "notes": "Weekend home",
  "from_roster": false
}
```

**Response (201 Created - Auto-verified - Apartment):**
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "user": {
      "user_id": "uuid-v4",
      "phone": "+919876543210",
      "full_name": "Ramesh Kumar",
      "user_type": "RESIDENT",
      "is_verified": true,
      "email": "ramesh@example.com"
    },
    "residences": [
      {
        "resident_id": "uuid-v4",
        "society_id": 1,
        "society_name": "Maple Gardens",
        "society_type": "APARTMENT",
        "unit_type": "FLAT",
        "flat_number": "A-404",
        "tower": "A",
        "floor": 4,
        "verification_status": "VERIFIED",
        "is_primary": true
      }
    ],
    "active_society": {
      "society_id": 1,
      "society_name": "Maple Gardens"
    },
    "tokens": {
      "access_token": "jwt-token",
      "refresh_token": "jwt-refresh-token",
      "expires_in": 3600
    }
  }
}
```

**Response (201 Created - Auto-verified - Independent House):**
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "user": {
      "user_id": "uuid-v4",
      "phone": "+919876543211",
      "full_name": "Priya Sharma",
      "user_type": "RESIDENT",
      "is_verified": true,
      "email": "priya@example.com"
    },
    "residences": [
      {
        "resident_id": "uuid-v4",
        "society_id": 5,
        "society_name": "Green Meadows Layout",
        "society_type": "LAYOUT",
        "unit_type": "HOUSE",
        "house_number": "42",
        "street": "5th Cross",
        "floor": 1,
        "notes": "First floor",
        "verification_status": "VERIFIED",
        "is_primary": true
      }
    ],
    "active_society": {
      "society_id": 5,
      "society_name": "Green Meadows Layout"
    },
    "tokens": {
      "access_token": "jwt-token",
      "refresh_token": "jwt-refresh-token",
      "expires_in": 3600
    }
  }
}
```

**Response (201 Created - Pending Verification):**
```json
{
  "success": true,
  "message": "Registration submitted for verification",
  "data": {
    "user": {
      "user_id": "uuid-v4",
      "phone": "+919876543210",
      "full_name": "Ramesh Kumar",
      "user_type": "RESIDENT",
      "is_verified": false
    },
    "residences": [
      {
        "resident_id": "uuid-v4",
        "society_id": 1,
        "society_name": "Maple Gardens",
        "unit_type": "FLAT",
        "flat_number": "B-302",
        "tower": "B",
        "floor": 3,
        "verification_status": "PENDING",
        "is_primary": true
      }
    ],
    "active_society": {
      "society_id": 1,
      "society_name": "Maple Gardens"
    },
    "tokens": {
      "access_token": "jwt-token",
      "refresh_token": "jwt-refresh-token",
      "expires_in": 3600
    },
    "next_steps": {
      "message": "Your registration is pending society admin approval",
      "estimated_time": "24 hours",
      "can_browse": true,
      "can_place_orders": false
    }
  }
}
```

**Response (200 OK - Additional Society Added):**
```json
{
  "success": true,
  "message": "Additional society added successfully",
  "data": {
    "user_id": "existing-uuid-v4",
    "residences": [
      {
        "resident_id": "uuid-v4-existing",
        "society_id": 1,
        "society_name": "Maple Gardens",
        "unit_type": "FLAT",
        "flat_number": "A-404",
        "verification_status": "VERIFIED",
        "is_primary": true
      },
      {
        "resident_id": "uuid-v4-new",
        "society_id": 3,
        "society_name": "Beach View Apartments",
        "unit_type": "FLAT",
        "flat_number": "201",
        "floor": 2,
        "notes": "Weekend home",
        "verification_status": "PENDING",
        "is_primary": false
      }
    ],
    "active_society": {
      "society_id": 1,
      "society_name": "Maple Gardens"
    }
  }
}
```

**Notes:**
- `unit_type`: `FLAT` for apartments, `HOUSE` for independent houses
- For apartments: `flat_number`, `tower` (optional), `floor` are used
- For independent houses: `house_number`, `street` (optional), `floor` are used
- `floor`: 0 = ground floor, 1 = first floor, 2 = second floor, etc.
- Multiple households in same house: Same `house_number`, different `floor` values
- Multi-society users: Only one residence can be `is_primary: true`
- Additional societies always require verification (even if phone in roster)

---

#### 2.1.4 Switch Active Society

**Endpoint:** `POST /api/v1/residents/{user_id}/switch-society`

**Description:** Switch the active society context for a multi-society resident. Changes which society's vendors and services are shown.

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body:**
```json
{
  "society_id": 3
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Active society switched successfully",
  "data": {
    "user_id": "uuid-v4",
    "previous_society": {
      "society_id": 1,
      "society_name": "Maple Gardens"
    },
    "active_society": {
      "society_id": 3,
      "society_name": "Beach View Apartments",
      "society_type": "APARTMENT",
      "address": "456 Beach Road, Chennai",
      "city": "Chennai"
    },
    "residence": {
      "resident_id": "uuid-v4-2",
      "unit_type": "FLAT",
      "flat_number": "201",
      "floor": 2,
      "verification_status": "VERIFIED"
    },
    "available_vendors": 23,
    "active_categories": ["LAUNDRY"]
  }
}
```

**Response (403 Forbidden - Society not verified):**
```json
{
  "success": false,
  "error": {
    "code": "SOCIETY_NOT_VERIFIED",
    "message": "You are not verified as a resident of Beach View Apartments",
    "details": {
      "society_id": 3,
      "society_name": "Beach View Apartments",
      "verification_status": "PENDING",
      "submitted_at": "2025-11-19T10:00:00Z"
    }
  }
}
```

**Response (404 Not Found):**
```json
{
  "success": false,
  "error": {
    "code": "SOCIETY_NOT_FOUND",
    "message": "You are not registered in the specified society",
    "available_societies": [
      {
        "society_id": 1,
        "society_name": "Maple Gardens"
      }
    ]
  }
}
```

**Notes:**
- All subsequent API calls (vendor listing, order creation) use the active society context
- Active society preference is stored in user session/profile
- User can only switch to societies where they have a verified residence
- Switching society updates the vendor list, rate cards shown, and available services

---

#### 2.1.5 Get User Residences

**Endpoint:** `GET /api/v1/residents/{user_id}/residences`

**Description:** Get all residences/societies for a user

**Headers:**
```
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "user_id": "uuid-v4",
    "residences": [
      {
        "resident_id": "uuid-v4-1",
        "society_id": 1,
        "society_name": "Maple Gardens",
        "society_type": "APARTMENT",
        "address": "123 MG Road, Koramangala, Bangalore",
        "city": "Bangalore",
        "unit_type": "FLAT",
        "flat_number": "A-404",
        "tower": "A",
        "floor": 4,
        "verification_status": "VERIFIED",
        "is_primary": true,
        "is_active": true,
        "verified_at": "2025-11-01T10:00:00Z"
      },
      {
        "resident_id": "uuid-v4-2",
        "society_id": 3,
        "society_name": "Beach View Apartments",
        "society_type": "APARTMENT",
        "address": "456 Beach Road, Chennai",
        "city": "Chennai",
        "unit_type": "FLAT",
        "flat_number": "201",
        "floor": 2,
        "verification_status": "VERIFIED",
        "is_primary": false,
        "is_active": false,
        "notes": "Weekend home",
        "verified_at": "2025-11-05T14:00:00Z"
      },
      {
        "resident_id": "uuid-v4-3",
        "society_id": 7,
        "society_name": "Hill View Layout",
        "society_type": "LAYOUT",
        "address": "House #15, Hill View Road, Ooty",
        "city": "Ooty",
        "unit_type": "HOUSE",
        "house_number": "15",
        "floor": 2,
        "verification_status": "PENDING",
        "is_primary": false,
        "is_active": false,
        "notes": "First floor - Rented",
        "submitted_at": "2025-11-19T10:00:00Z"
      }
    ],
    "active_society": {
      "society_id": 1,
      "society_name": "Maple Gardens"
    },
    "total_residences": 3,
    "verified_count": 2,
    "pending_count": 1
  }
}
```

**Notes:**
- `is_active`: Indicates which society is currently selected/active
- `is_primary`: User's main residence (set during initial registration)
- Users can have residences in pending verification status
- Pending residences cannot be set as active until verified

---

### 2.2 Vendor Onboarding

#### 2.2.1 Initiate Vendor Registration

**Endpoint:** `POST /api/v1/onboarding/vendor/register`

**Description:** Start vendor registration process

**Request Body:**
```json
{
  "phone": "+919876543211",
  "full_name": "Priya Sharma",
  "email": "priya@perfectpress.com",
  "business_name": "Perfect Press",
  "store_address": "789 Market Street, Koramangala",
  "id_proof_type": "AADHAAR",
  "id_proof_number": "1234-5678-9012",
  "id_proof_photo_url": "https://...",
  "store_photo_url": "https://...",
  "gst_number": "29ABCDE1234F1Z5",
  "pan_number": "ABCDE1234F"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Vendor registration submitted successfully",
  "data": {
    "user": {
      "user_id": "uuid-v4",
      "phone": "+919876543211",
      "full_name": "Priya Sharma",
      "user_type": "VENDOR",
      "is_verified": false
    },
    "vendor": {
      "vendor_id": "uuid-v4",
      "business_name": "Perfect Press",
      "store_address": "789 Market Street, Koramangala",
      "approval_status": "PENDING",
      "created_at": "2025-11-20T10:00:00Z"
    },
    "tokens": {
      "access_token": "jwt-token",
      "refresh_token": "jwt-refresh-token",
      "expires_in": 3600
    },
    "next_steps": {
      "message": "Registration submitted. Complete your profile to get verified",
      "required_actions": [
        "Add bank account details",
        "Select services offered",
        "Choose societies to serve",
        "Create rate cards"
      ]
    }
  }
}
```

---

#### 2.2.2 Update Vendor Bank Details

**Endpoint:** `PUT /api/v1/onboarding/vendor/{vendor_id}/bank-details`

**Description:** Add or update vendor bank account information

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body:**
```json
{
  "bank_account_number": "1234567890123",
  "bank_ifsc_code": "SBIN0001234",
  "bank_account_holder": "Priya Sharma",
  "bank_name": "State Bank of India",
  "branch_name": "Koramangala Branch"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Bank details updated successfully",
  "data": {
    "vendor_id": "uuid-v4",
    "bank_details": {
      "bank_account_number": "***********3",
      "bank_ifsc_code": "SBIN0001234",
      "bank_account_holder": "Priya Sharma",
      "bank_name": "State Bank of India",
      "verified": false
    }
  }
}
```

---

#### 2.2.3 Select Services Offered

**Endpoint:** `POST /api/v1/onboarding/vendor/{vendor_id}/services`

**Description:** Select which service categories and types the vendor offers

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body:**
```json
{
  "services": [
    {
      "service_id": 1,
      "service_key": "IRONING",
      "turnaround_hours": 24,
      "is_active": true
    },
    {
      "service_id": 2,
      "service_key": "WASHING_IRONING",
      "turnaround_hours": 48,
      "is_active": true
    },
    {
      "service_id": 3,
      "service_key": "DRY_CLEANING",
      "turnaround_hours": 120,
      "is_active": true
    }
  ]
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Services updated successfully",
  "data": {
    "vendor_id": "uuid-v4",
    "services_offered": [
      {
        "service_id": 1,
        "service_name": "Ironing Only",
        "service_key": "IRONING",
        "category": "LAUNDRY",
        "turnaround_hours": 24,
        "is_active": true
      },
      {
        "service_id": 2,
        "service_name": "Washing + Ironing",
        "service_key": "WASHING_IRONING",
        "category": "LAUNDRY",
        "turnaround_hours": 48,
        "is_active": true
      },
      {
        "service_id": 3,
        "service_name": "Dry Cleaning",
        "service_key": "DRY_CLEANING",
        "category": "LAUNDRY",
        "turnaround_hours": 120,
        "is_active": true
      }
    ]
  }
}
```

---

#### 2.2.4 Request Society Access

**Endpoint:** `POST /api/v1/onboarding/vendor/{vendor_id}/societies`

**Description:** Request access to serve specific societies

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body:**
```json
{
  "society_ids": [1, 2, 3]
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Society access requests submitted",
  "data": {
    "vendor_id": "uuid-v4",
    "society_requests": [
      {
        "society_id": 1,
        "society_name": "Maple Gardens",
        "approval_status": "PENDING",
        "requested_at": "2025-11-20T10:00:00Z",
        "estimated_approval_time": "24-48 hours"
      },
      {
        "society_id": 2,
        "society_name": "Palm Residency",
        "approval_status": "PENDING",
        "requested_at": "2025-11-20T10:00:00Z",
        "estimated_approval_time": "24-48 hours"
      }
    ]
  }
}
```

---

## 3. Approval & Verification APIs

### 3.1 Society Admin Approvals

#### 3.1.1 Get Pending Resident Approvals

**Endpoint:** `GET /api/v1/admin/society/{society_id}/residents/pending`

**Description:** Get list of residents pending approval for a society

**Headers:**
```
Authorization: Bearer {access_token}
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20)
- `sort_by` (optional): created_at, flat_number (default: created_at)
- `order` (optional): asc, desc (default: desc)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "residents": [
      {
        "resident_id": "uuid-v4",
        "user_id": "uuid-v4",
        "full_name": "Ramesh Kumar",
        "phone": "+919876543210",
        "email": "ramesh@example.com",
        "flat_number": "B-302",
        "tower": "B",
        "floor": 3,
        "verification_status": "PENDING",
        "created_at": "2025-11-20T09:00:00Z",
        "days_pending": 0
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 1,
      "total_items": 1,
      "items_per_page": 20
    }
  }
}
```

---

#### 3.1.2 Approve/Reject Resident

**Endpoint:** `POST /api/v1/admin/society/{society_id}/residents/{resident_id}/approve`

**Description:** Approve or reject a resident registration

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body (Approve):**
```json
{
  "action": "APPROVE",
  "notes": "Verified from society records"
}
```

**Request Body (Reject):**
```json
{
  "action": "REJECT",
  "rejection_reason": "Flat number does not exist in society records",
  "notes": "Please contact society office for correct flat number"
}
```

**Response (200 OK - Approved):**
```json
{
  "success": true,
  "message": "Resident approved successfully",
  "data": {
    "resident_id": "uuid-v4",
    "full_name": "Ramesh Kumar",
    "flat_number": "B-302",
    "verification_status": "VERIFIED",
    "approved_at": "2025-11-20T10:00:00Z",
    "approved_by": "uuid-admin-id",
    "notification_sent": true
  }
}
```

**Response (200 OK - Rejected):**
```json
{
  "success": true,
  "message": "Resident registration rejected",
  "data": {
    "resident_id": "uuid-v4",
    "full_name": "Ramesh Kumar",
    "verification_status": "REJECTED",
    "rejection_reason": "Flat number does not exist in society records",
    "rejected_at": "2025-11-20T10:00:00Z",
    "rejected_by": "uuid-admin-id",
    "notification_sent": true
  }
}
```

---

#### 3.1.3 Get Pending Vendor Approvals

**Endpoint:** `GET /api/v1/admin/society/{society_id}/vendors/pending`

**Description:** Get list of vendors requesting access to the society

**Headers:**
```
Authorization: Bearer {access_token}
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "vendors": [
      {
        "vendor_id": "uuid-v4",
        "business_name": "Perfect Press",
        "full_name": "Priya Sharma",
        "phone": "+919876543211",
        "store_address": "789 Market Street, Koramangala",
        "services_offered": [
          {
            "service_id": 1,
            "service_name": "Ironing Only",
            "turnaround_hours": 24
          },
          {
            "service_id": 2,
            "service_name": "Washing + Ironing",
            "turnaround_hours": 48
          }
        ],
        "approval_status": "PENDING",
        "requested_at": "2025-11-20T09:00:00Z",
        "has_rate_card": false,
        "platform_approval_status": "PENDING",
        "total_orders": 0,
        "avg_rating": 0
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 1,
      "total_items": 1,
      "items_per_page": 20
    }
  }
}
```

---

#### 3.1.4 Approve/Reject Vendor for Society

**Endpoint:** `POST /api/v1/admin/society/{society_id}/vendors/{vendor_id}/approve`

**Description:** Approve or reject vendor access to the society

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body (Approve):**
```json
{
  "action": "APPROVE",
  "notes": "Good reviews from other societies"
}
```

**Request Body (Reject):**
```json
{
  "action": "REJECT",
  "rejection_reason": "Society already has sufficient laundry vendors",
  "notes": "May reconsider in future"
}
```

**Response (200 OK - Approved):**
```json
{
  "success": true,
  "message": "Vendor approved for society",
  "data": {
    "vendor_id": "uuid-v4",
    "business_name": "Perfect Press",
    "society_id": 1,
    "approval_status": "APPROVED",
    "approved_at": "2025-11-20T10:00:00Z",
    "approved_by": "uuid-admin-id",
    "notification_sent": true,
    "next_steps": {
      "message": "Vendor must create rate card for this society before going live",
      "has_rate_card": false
    }
  }
}
```

---

#### 3.1.5 Manage Society Groups (Unified for Buildings/Phases)

**Endpoint:** `POST /api/v1/admin/society/{society_id}/groups`

**Description:** Create a new group (building/phase/tower/block/etc.) - unified for both apartment and layout societies

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body (Apartment - Building):**
```json
{
  "group_name": "Building A",
  "group_code": "A",
  "group_type": "BUILDING",
  "description": "Main residential tower",
  "total_floors": 15,
  "total_units": 60
}
```

**Request Body (Layout - Phase):**
```json
{
  "group_name": "Phase 1",
  "group_code": "P1",
  "group_type": "PHASE",
  "description": "Eastern section of the layout",
  "total_units": 50
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Group created successfully",
  "data": {
    "group_id": 1,
    "society_id": 1,
    "group_name": "Building A",
    "group_code": "A",
    "group_type": "BUILDING",
    "total_floors": 15,
    "total_units": 60,
    "created_at": "2025-11-20T10:00:00Z"
  }
}
```

**Notes:**
- **Unified endpoint** for both apartments and layouts
- `group_type` options: 'BUILDING', 'BLOCK', 'TOWER', 'WING', 'PHASE', 'SECTION', 'ZONE'
- `total_units`: Number of flats (for apartments) OR houses (for layouts)
- `total_floors`: Only applicable for multi-story buildings

---

**Endpoint:** `GET /api/v1/admin/society/{society_id}/groups`

**Description:** Get all groups for a society (works for both apartments and layouts)

**Query Parameters:**
- `group_type` (optional): Filter by type (BUILDING, PHASE, etc.)

**Response (200 OK - Apartment Society):**
```json
{
  "success": true,
  "data": {
    "society_id": 1,
    "society_type": "APARTMENT",
    "groups": [
      {
        "group_id": 1,
        "group_name": "Building A",
        "group_code": "A",
        "group_type": "BUILDING",
        "total_floors": 15,
        "total_units": 60,
        "is_active": true
      },
      {
        "group_id": 2,
        "group_name": "Tower B",
        "group_code": "B",
        "group_type": "TOWER",
        "total_floors": 20,
        "total_units": 80,
        "is_active": true
      }
    ]
  }
}
```

**Response (200 OK - Layout Society):**
```json
{
  "success": true,
  "data": {
    "society_id": 2,
    "society_type": "LAYOUT",
    "groups": [
      {
        "group_id": 5,
        "group_name": "Phase 1",
        "group_code": "P1",
        "group_type": "PHASE",
        "total_units": 50,
        "is_active": true
      },
      {
        "group_id": 6,
        "group_name": "East Section",
        "group_code": "ES",
        "group_type": "SECTION",
        "total_units": 35,
        "is_active": true
      }
    ]
  }
}
```

---

#### 3.1.6 Assign Vendor to Service Areas (Unified Groups)

**Endpoint:** `POST /api/v1/admin/society/{society_id}/vendors/{vendor_id}/service-areas`

**Description:** Assign vendor to specific groups (buildings/phases) or entire society - unified approach

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body (Assign to entire society):**
```json
{
  "assignment_type": "SOCIETY"
}
```

**Request Body (Assign to specific groups - works for both buildings and phases):**
```json
{
  "assignment_type": "GROUP",
  "group_ids": [1, 2]
}
```

**Response (200 OK - Buildings):**
```json
{
  "success": true,
  "message": "Vendor service areas assigned successfully",
  "data": {
    "vendor_id": "uuid-v4",
    "business_name": "Perfect Press",
    "society_id": 1,
    "assignments": [
      {
        "assignment_id": 1,
        "assignment_type": "GROUP",
        "group_id": 1,
        "group_name": "Building A",
        "group_type": "BUILDING",
        "is_active": true
      },
      {
        "assignment_id": 2,
        "assignment_type": "GROUP",
        "group_id": 2,
        "group_name": "Building B",
        "group_type": "BUILDING",
        "is_active": true
      }
    ],
    "coverage_summary": {
      "covers_entire_society": false,
      "groups_assigned": 2,
      "total_groups_in_society": 5,
      "estimated_households": 108
    }
  }
}
```

**Response (200 OK - Phases):**
```json
{
  "success": true,
  "message": "Vendor service areas assigned successfully",
  "data": {
    "vendor_id": "uuid-v4",
    "business_name": "Express Cleaners",
    "society_id": 2,
    "assignments": [
      {
        "assignment_id": 5,
        "assignment_type": "GROUP",
        "group_id": 5,
        "group_name": "Phase 1",
        "group_type": "PHASE",
        "is_active": true
      },
      {
        "assignment_id": 6,
        "assignment_type": "GROUP",
        "group_id": 6,
        "group_name": "Phase 2",
        "group_type": "PHASE",
        "is_active": true
      }
    ],
    "coverage_summary": {
      "covers_entire_society": false,
      "groups_assigned": 2,
      "total_groups_in_society": 3,
      "estimated_households": 85
    }
  }
}
```

**Notes:**
- **Simplified assignment types:** Only 'SOCIETY' or 'GROUP'
- Works uniformly for both apartment buildings and layout phases
- Can mix different group types in the same assignment (e.g., Building A + Tower B)

---

#### 3.1.7 Get Vendor Service Area Assignments

**Endpoint:** `GET /api/v1/admin/society/{society_id}/vendors/{vendor_id}/service-areas`

**Description:** Get current service area assignments for a vendor

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "vendor_id": "uuid-v4",
    "business_name": "Perfect Press",
    "society_id": 1,
    "assignments": [
      {
        "assignment_id": 1,
        "assignment_type": "GROUP",
        "group_id": 1,
        "group_name": "Building A",
        "group_type": "BUILDING",
        "is_active": true,
        "assigned_at": "2025-11-20T10:00:00Z"
      },
      {
        "assignment_id": 2,
        "assignment_type": "GROUP",
        "group_id": 2,
        "group_name": "Tower B",
        "group_type": "TOWER",
        "is_active": true,
        "assigned_at": "2025-11-20T10:00:00Z"
      }
    ],
    "coverage_summary": {
      "covers_entire_society": false,
      "assignment_level": "GROUP",
      "groups_assigned": ["Building A", "Tower B"],
      "total_groups_in_society": 5,
      "coverage_percentage": 40
    }
  }
}
```

---

#### 3.1.8 Update Vendor Service Area Assignments

**Endpoint:** `PUT /api/v1/admin/society/{society_id}/vendors/{vendor_id}/service-areas`

**Description:** Update vendor's service area assignments (replaces existing assignments)

**Request Body (Change to society-wide):**
```json
{
  "assignment_type": "SOCIETY"
}
```

**Request Body (Change to specific groups):**
```json
{
  "assignment_type": "GROUP",
  "group_ids": [1, 2, 3]
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Vendor service areas updated successfully",
  "data": {
    "vendor_id": "uuid-v4",
    "previous_assignments": [
      {
        "assignment_type": "GROUP",
        "groups": ["Building A", "Building B"]
      }
    ],
    "new_assignments": [
      {
        "assignment_type": "SOCIETY",
        "covers_entire_society": true
      }
    ],
    "updated_at": "2025-11-20T11:00:00Z"
  }
}
```

---

#### 3.1.9 Delete Vendor Service Area Assignment

**Endpoint:** `DELETE /api/v1/admin/society/{society_id}/vendors/{vendor_id}/service-areas/{assignment_id}`

**Description:** Remove a specific service area assignment

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Service area assignment removed successfully",
  "data": {
    "assignment_id": 1,
    "vendor_id": "uuid-v4",
    "remaining_assignments": 1
  }
}
```

---

### 3.2 Platform Admin Approvals

#### 3.2.1 Get Pending Vendor Verifications

**Endpoint:** `GET /api/v1/admin/platform/vendors/pending`

**Description:** Get all vendors pending platform-level verification

**Headers:**
```
Authorization: Bearer {access_token}
```

**Query Parameters:**
- `page` (optional): Page number
- `limit` (optional): Items per page

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "vendors": [
      {
        "vendor_id": "uuid-v4",
        "business_name": "Perfect Press",
        "full_name": "Priya Sharma",
        "phone": "+919876543211",
        "email": "priya@perfectpress.com",
        "store_address": "789 Market Street, Koramangala",
        "id_proof_type": "AADHAAR",
        "id_proof_number": "1234-5678-9012",
        "id_proof_photo_url": "https://...",
        "store_photo_url": "https://...",
        "gst_number": "29ABCDE1234F1Z5",
        "pan_number": "ABCDE1234F",
        "bank_account_number": "***********3",
        "bank_ifsc_code": "SBIN0001234",
        "services_offered": ["IRONING", "WASHING_IRONING", "DRY_CLEANING"],
        "societies_requested": 3,
        "approval_status": "PENDING",
        "created_at": "2025-11-20T09:00:00Z",
        "days_pending": 0
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 1,
      "total_items": 1
    }
  }
}
```

---

#### 3.2.2 Approve/Reject Vendor (Platform Level)

**Endpoint:** `POST /api/v1/admin/platform/vendors/{vendor_id}/verify`

**Description:** Platform admin verifies vendor identity and business details

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body (Approve):**
```json
{
  "action": "APPROVE",
  "verification_notes": "All documents verified",
  "verified_fields": {
    "id_proof": true,
    "business_registration": true,
    "bank_details": true,
    "gst_pan": true
  }
}
```

**Request Body (Reject):**
```json
{
  "action": "REJECT",
  "rejection_reason": "Invalid GST number",
  "notes": "Please provide valid GST registration certificate"
}
```

**Response (200 OK - Approved):**
```json
{
  "success": true,
  "message": "Vendor verified successfully",
  "data": {
    "vendor_id": "uuid-v4",
    "business_name": "Perfect Press",
    "approval_status": "APPROVED",
    "is_verified": true,
    "approved_at": "2025-11-20T10:00:00Z",
    "approved_by": "uuid-admin-id",
    "notification_sent": true,
    "next_steps": {
      "message": "Vendor can now request access to societies and create rate cards"
    }
  }
}
```

---

## 4. Rate Card Management APIs

### 4.1 Create Rate Card for Society

**Endpoint:** `POST /api/v1/vendors/{vendor_id}/rate-cards`

**Description:** Create a new rate card for a specific society

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body:**
```json
{
  "society_id": 1,
  "items": [
    {
      "service_id": 1,
      "service_key": "IRONING",
      "items": [
        {
          "item_name": "Shirt",
          "description": "Regular shirt ironing",
          "price_per_piece": 10.00,
          "display_order": 1
        },
        {
          "item_name": "Pants",
          "description": "Regular pants ironing",
          "price_per_piece": 15.00,
          "display_order": 2
        },
        {
          "item_name": "Saree",
          "description": "Saree ironing",
          "price_per_piece": 30.00,
          "display_order": 3
        }
      ]
    },
    {
      "service_id": 2,
      "service_key": "WASHING_IRONING",
      "items": [
        {
          "item_name": "Shirt",
          "description": "Wash and iron",
          "price_per_piece": 25.00,
          "display_order": 1
        },
        {
          "item_name": "Pants",
          "description": "Wash and iron",
          "price_per_piece": 30.00,
          "display_order": 2
        }
      ]
    },
    {
      "service_id": 3,
      "service_key": "DRY_CLEANING",
      "items": [
        {
          "item_name": "Shirt",
          "description": "Professional dry cleaning",
          "price_per_piece": 80.00,
          "display_order": 1
        },
        {
          "item_name": "Blazer",
          "description": "Professional dry cleaning",
          "price_per_piece": 150.00,
          "display_order": 2
        },
        {
          "item_name": "Saree",
          "description": "Professional dry cleaning",
          "price_per_piece": 200.00,
          "display_order": 3
        }
      ]
    }
  ]
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Rate card created successfully",
  "data": {
    "rate_card_id": 1,
    "vendor_id": "uuid-v4",
    "society_id": 1,
    "society_name": "Maple Gardens",
    "is_active": true,
    "is_published": false,
    "total_items": 8,
    "services_covered": [
      {
        "service_id": 1,
        "service_name": "Ironing Only",
        "items_count": 3,
        "price_range": {
          "min": 10.00,
          "max": 30.00
        }
      },
      {
        "service_id": 2,
        "service_name": "Washing + Ironing",
        "items_count": 2,
        "price_range": {
          "min": 25.00,
          "max": 30.00
        }
      },
      {
        "service_id": 3,
        "service_name": "Dry Cleaning",
        "items_count": 3,
        "price_range": {
          "min": 80.00,
          "max": 200.00
        }
      }
    ],
    "created_at": "2025-11-20T10:00:00Z",
    "next_steps": {
      "message": "Review and publish rate card to make it visible to residents"
    }
  }
}
```

---

### 4.2 Get Rate Card for Vendor-Society

**Endpoint:** `GET /api/v1/vendors/{vendor_id}/rate-cards/{society_id}`

**Description:** Get rate card details for a specific vendor-society combination

**Headers:**
```
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "rate_card_id": 1,
    "vendor_id": "uuid-v4",
    "vendor_name": "Perfect Press",
    "society_id": 1,
    "society_name": "Maple Gardens",
    "is_active": true,
    "is_published": true,
    "published_at": "2025-11-20T11:00:00Z",
    "services": [
      {
        "service_id": 1,
        "service_name": "Ironing Only",
        "service_key": "IRONING",
        "items": [
          {
            "item_id": 1,
            "item_name": "Shirt",
            "description": "Regular shirt ironing",
            "price_per_piece": 10.00,
            "display_order": 1,
            "is_active": true
          },
          {
            "item_id": 2,
            "item_name": "Pants",
            "description": "Regular pants ironing",
            "price_per_piece": 15.00,
            "display_order": 2,
            "is_active": true
          }
        ]
      },
      {
        "service_id": 2,
        "service_name": "Washing + Ironing",
        "service_key": "WASHING_IRONING",
        "items": [
          {
            "item_id": 3,
            "item_name": "Shirt",
            "description": "Wash and iron",
            "price_per_piece": 25.00,
            "display_order": 1,
            "is_active": true
          }
        ]
      }
    ],
    "created_at": "2025-11-20T10:00:00Z",
    "updated_at": "2025-11-20T10:00:00Z"
  }
}
```

---

### 4.3 Update Rate Card Items

**Endpoint:** `PUT /api/v1/vendors/{vendor_id}/rate-cards/{rate_card_id}/items`

**Description:** Update existing rate card items or add new items

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body:**
```json
{
  "items": [
    {
      "item_id": 1,
      "price_per_piece": 12.00,
      "is_active": true
    },
    {
      "service_id": 1,
      "item_name": "Bedsheet",
      "description": "Single bedsheet ironing",
      "price_per_piece": 20.00,
      "display_order": 4
    }
  ]
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Rate card updated successfully",
  "data": {
    "rate_card_id": 1,
    "updated_items": 1,
    "added_items": 1,
    "total_items": 9,
    "is_published": false,
    "requires_republish": true,
    "message": "Changes saved as draft. Publish to make visible to residents"
  }
}
```

---

### 4.4 Publish Rate Card

**Endpoint:** `POST /api/v1/vendors/{vendor_id}/rate-cards/{rate_card_id}/publish`

**Description:** Publish rate card to make it visible to residents

**Headers:**
```
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Rate card published successfully",
  "data": {
    "rate_card_id": 1,
    "is_published": true,
    "published_at": "2025-11-20T11:00:00Z",
    "total_items": 9,
    "services_covered": 3,
    "message": "Your rate card is now visible to all residents in Maple Gardens"
  }
}
```

---

## 5. Vendor Listing & Discovery APIs

### 5.1 Get Service Categories

**Endpoint:** `GET /api/v1/categories`

**Description:** Get all parent service categories (Laundry, Vehicle, Home, Personal)

**Query Parameters:**
- `is_live` (optional): Filter by live status (true/false)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "categories": [
      {
        "category_id": 1,
        "category_key": "LAUNDRY",
        "category_name": "Laundry Services",
        "description": "Professional laundry and garment care services",
        "icon_url": "https://...",
        "color_hex": "#3B82F6",
        "is_live": true,
        "display_order": 1,
        "subcategories": [
          {
            "service_id": 1,
            "service_key": "IRONING",
            "service_name": "Ironing Only",
            "default_turnaround_hours": 24
          },
          {
            "service_id": 2,
            "service_key": "WASHING_IRONING",
            "service_name": "Washing + Ironing",
            "default_turnaround_hours": 48
          },
          {
            "service_id": 3,
            "service_key": "DRY_CLEANING",
            "service_name": "Dry Cleaning",
            "default_turnaround_hours": 120
          }
        ]
      },
      {
        "category_id": 2,
        "category_key": "VEHICLE",
        "category_name": "Vehicle Services",
        "description": "Car and bike washing and detailing services",
        "is_live": false,
        "display_order": 2,
        "message": "Coming Soon"
      }
    ]
  }
}
```

---

### 5.2 List Vendors for Society

**Endpoint:** `GET /api/v1/societies/{society_id}/vendors`

**Description:** Get approved vendors serving a specific society with smart filtering based on resident's location

**Headers:**
```
Authorization: Bearer {access_token}
```

**Query Parameters:**
- `category` (optional): Filter by category (LAUNDRY, VEHICLE, etc.)
- `service_id` (optional): Filter by specific service type
- `group_id` (optional): Filter by specific group (building/phase)
- `show_all` (optional): Show all vendors regardless of assignment (default: false)
- `sort_by` (optional): rating, delivery_time, orders (default: rating)
- `order` (optional): asc, desc (default: desc)
- `page` (optional): Page number
- `limit` (optional): Items per page

**Filtering Logic:**
- **Default behavior (show_all=false):** Returns vendors assigned to:
  - The resident's specific group (building/phase), OR
  - The entire society (assignment_type='SOCIETY')
- **Override behavior (show_all=true):** Returns ALL vendors serving the society regardless of assignment
- If `group_id` is explicitly provided, uses that for filtering instead of resident's group

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "vendors": [
      {
        "vendor_id": "uuid-v4",
        "business_name": "Perfect Press",
        "store_address": "789 Market Street, Koramangala",
        "store_photo_url": "https://...",
        "avg_rating": 4.5,
        "total_orders": 150,
        "completed_orders": 145,
        "is_available": true,
        "service_areas": {
          "assignment_type": "GROUP",
          "covers_entire_society": false,
          "assigned_groups": [
            {
              "group_id": 1,
              "group_name": "Building A",
              "group_type": "BUILDING"
            },
            {
              "group_id": 2,
              "group_name": "Building B",
              "group_type": "BUILDING"
            }
          ]
        },
        "services_offered": [
          {
            "service_id": 1,
            "service_name": "Ironing Only",
            "service_key": "IRONING",
            "category": "LAUNDRY",
            "turnaround_hours": 24,
            "starting_price": 10.00
          },
          {
            "service_id": 2,
            "service_name": "Washing + Ironing",
            "service_key": "WASHING_IRONING",
            "category": "LAUNDRY",
            "turnaround_hours": 48,
            "starting_price": 25.00
          },
          {
            "service_id": 3,
            "service_name": "Dry Cleaning",
            "service_key": "DRY_CLEANING",
            "category": "LAUNDRY",
            "turnaround_hours": 120,
            "starting_price": 80.00
          }
        ],
        "has_rate_card": true,
        "is_published": true
      },
      {
        "vendor_id": "uuid-v5",
        "business_name": "Express Cleaners",
        "store_address": "456 Service Road, Koramangala",
        "store_photo_url": "https://...",
        "avg_rating": 4.8,
        "total_orders": 280,
        "completed_orders": 275,
        "is_available": true,
        "service_areas": {
          "assignment_type": "SOCIETY",
          "covers_entire_society": true,
          "assigned_groups": null
        },
        "services_offered": [
          {
            "service_id": 1,
            "service_name": "Ironing Only",
            "service_key": "IRONING",
            "category": "LAUNDRY",
            "turnaround_hours": 24,
            "starting_price": 12.00
          }
        ],
        "has_rate_card": true,
        "is_published": true
      }
    ],
    "filter_info": {
      "show_all": false,
      "filtered_by_group": true,
      "group_id": 1,
      "group_name": "Building A",
      "group_type": "BUILDING",
      "total_vendors_in_society": 15,
      "vendors_shown": 2
    },
    "pagination": {
      "current_page": 1,
      "total_pages": 1,
      "total_items": 2
    }
  }
}
```

**Example Request - Default filtering for resident in Building A:**
```
GET /api/v1/societies/1/vendors?category=LAUNDRY
```
Returns vendors assigned to Building A + vendors assigned to entire society

**Example Request - View all vendors:**
```
GET /api/v1/societies/1/vendors?category=LAUNDRY&show_all=true
```
Returns ALL vendors in the society regardless of building/phase assignment

**Example Request - Specific group:**
```
GET /api/v1/societies/1/vendors?group_id=2
```
Returns vendors assigned to group ID 2 (e.g., Building B) + vendors assigned to entire society

---

### 5.3 Get Vendor Details

**Endpoint:** `GET /api/v1/vendors/{vendor_id}`

**Description:** Get detailed vendor information including services and rate card

**Headers:**
```
Authorization: Bearer {access_token}
```

**Query Parameters:**
- `society_id` (required): Society context for rate card

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "vendor": {
      "vendor_id": "uuid-v4",
      "business_name": "Perfect Press",
      "full_name": "Priya Sharma",
      "phone": "+919876543211",
      "store_address": "789 Market Street, Koramangala",
      "store_photo_url": "https://...",
      "avg_rating": 4.5,
      "total_orders": 150,
      "completed_orders": 145,
      "is_available": true,
      "services_offered": [
        {
          "service_id": 1,
          "service_name": "Ironing Only",
          "service_key": "IRONING",
          "category": "LAUNDRY",
          "turnaround_hours": 24,
          "pricing_model": "PER_ITEM"
        },
        {
          "service_id": 2,
          "service_name": "Washing + Ironing",
          "service_key": "WASHING_IRONING",
          "category": "LAUNDRY",
          "turnaround_hours": 48,
          "pricing_model": "PER_ITEM"
        },
        {
          "service_id": 3,
          "service_name": "Dry Cleaning",
          "service_key": "DRY_CLEANING",
          "category": "LAUNDRY",
          "turnaround_hours": 120,
          "pricing_model": "PER_ITEM"
        }
      ],
      "rate_card": {
        "rate_card_id": 1,
        "society_id": 1,
        "is_published": true,
        "services": [
          {
            "service_id": 1,
            "service_name": "Ironing Only",
            "items": [
              {
                "item_id": 1,
                "item_name": "Shirt",
                "price_per_piece": 10.00
              },
              {
                "item_id": 2,
                "item_name": "Pants",
                "price_per_piece": 15.00
              },
              {
                "item_id": 3,
                "item_name": "Saree",
                "price_per_piece": 30.00
              }
            ]
          },
          {
            "service_id": 2,
            "service_name": "Washing + Ironing",
            "items": [
              {
                "item_id": 4,
                "item_name": "Shirt",
                "price_per_piece": 25.00
              },
              {
                "item_id": 5,
                "item_name": "Pants",
                "price_per_piece": 30.00
              }
            ]
          },
          {
            "service_id": 3,
            "service_name": "Dry Cleaning",
            "items": [
              {
                "item_id": 6,
                "item_name": "Shirt",
                "price_per_piece": 80.00
              },
              {
                "item_id": 7,
                "item_name": "Blazer",
                "price_per_piece": 150.00
              },
              {
                "item_id": 8,
                "item_name": "Saree",
                "price_per_piece": 200.00
              }
            ]
          }
        ]
      },
      "recent_reviews": [
        {
          "rating_id": 1,
          "rating": 5,
          "review": "Excellent service! Clothes were perfectly ironed",
          "service_id": 1,
          "service_name": "Ironing Only",
          "resident_name": "Ramesh K.",
          "created_at": "2025-11-19T15:00:00Z"
        }
      ]
    }
  }
}
```

---

## 6. Order Management APIs

### 6.1 Create Order

**Endpoint:** `POST /api/v1/orders`

**Description:** Create a new order with multiple service types

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body:**
```json
{
  "vendor_id": "uuid-v4",
  "society_id": 1,
  "pickup_datetime": "2025-11-21T10:30:00Z",
  "pickup_address": "A-404, Tower A, Maple Gardens",
  "resident_notes": "Please call before arriving",
  "items": [
    {
      "service_id": 1,
      "rate_card_item_id": 1,
      "item_name": "Shirt",
      "quantity": 5,
      "unit_price": 10.00
    },
    {
      "service_id": 2,
      "rate_card_item_id": 4,
      "item_name": "Shirt",
      "quantity": 3,
      "unit_price": 25.00
    },
    {
      "service_id": 2,
      "rate_card_item_id": 5,
      "item_name": "Pants",
      "quantity": 2,
      "unit_price": 30.00
    },
    {
      "service_id": 3,
      "rate_card_item_id": 7,
      "item_name": "Blazer",
      "quantity": 1,
      "unit_price": 150.00
    }
  ]
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "order": {
      "order_id": "uuid-v4",
      "order_number": "ORD20251120000001",
      "status": "PICKUP_SCHEDULED",
      "vendor": {
        "vendor_id": "uuid-v4",
        "business_name": "Perfect Press",
        "phone": "+919876543211"
      },
      "resident": {
        "resident_id": "uuid-v4",
        "full_name": "Ramesh Kumar",
        "phone": "+919876543210",
        "flat_number": "A-404"
      },
      "pickup_datetime": "2025-11-21T10:30:00Z",
      "pickup_address": "A-404, Tower A, Maple Gardens",
      "has_multiple_services": true,
      "services_summary": [
        {
          "service_id": 1,
          "service_name": "Ironing Only",
          "item_count": 5,
          "total_amount": 50.00,
          "expected_delivery_days": 1
        },
        {
          "service_id": 2,
          "service_name": "Washing + Ironing",
          "item_count": 5,
          "total_amount": 135.00,
          "expected_delivery_days": 2
        },
        {
          "service_id": 3,
          "service_name": "Dry Cleaning",
          "item_count": 1,
          "total_amount": 150.00,
          "expected_delivery_days": 5
        }
      ],
      "pricing": {
        "estimated_item_count": 11,
        "estimated_price": 335.00,
        "discount_amount": 0.00,
        "final_price": null
      },
      "expected_delivery_date": "2025-11-26",
      "delivery_note": "Based on Dry Cleaning turnaround (5 days)",
      "created_at": "2025-11-20T10:00:00Z"
    },
    "notifications": {
      "vendor_notified": true,
      "resident_confirmation_sent": true,
      "reminder_scheduled": "2025-11-21T10:00:00Z"
    }
  }
}
```

---

### 6.2 Get Order Details

**Endpoint:** `GET /api/v1/orders/{order_id}`

**Description:** Get complete order details including workflow progress

**Headers:**
```
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "order": {
      "order_id": "uuid-v4",
      "order_number": "ORD20251120000001",
      "status": "PROCESSING_IN_PROGRESS",
      "vendor": {
        "vendor_id": "uuid-v4",
        "business_name": "Perfect Press",
        "phone": "+919876543211",
        "store_address": "789 Market Street, Koramangala"
      },
      "resident": {
        "resident_id": "uuid-v4",
        "full_name": "Ramesh Kumar",
        "phone": "+919876543210",
        "flat_number": "A-404",
        "society_name": "Maple Gardens"
      },
      "pickup_datetime": "2025-11-21T10:30:00Z",
      "pickup_address": "A-404, Tower A, Maple Gardens",
      "expected_delivery_date": "2025-11-26",
      "has_multiple_services": true,
      "items": [
        {
          "item_name": "Shirt",
          "service_id": 1,
          "service_name": "Ironing Only",
          "quantity": 5,
          "unit_price": 10.00,
          "total_price": 50.00
        },
        {
          "item_name": "Shirt",
          "service_id": 2,
          "service_name": "Washing + Ironing",
          "quantity": 3,
          "unit_price": 25.00,
          "total_price": 75.00
        },
        {
          "item_name": "Pants",
          "service_id": 2,
          "service_name": "Washing + Ironing",
          "quantity": 2,
          "unit_price": 30.00,
          "total_price": 60.00
        },
        {
          "item_name": "Blazer",
          "service_id": 3,
          "service_name": "Dry Cleaning",
          "quantity": 1,
          "unit_price": 150.00,
          "total_price": 150.00
        }
      ],
      "service_status": [
        {
          "service_id": 1,
          "service_name": "Ironing Only",
          "status": "READY_FOR_DELIVERY",
          "item_count": 5,
          "total_amount": 50.00,
          "current_step": "Quality Check",
          "current_step_order": 4,
          "expected_delivery_date": "2025-11-22",
          "ready_at": "2025-11-21T18:00:00Z"
        },
        {
          "service_id": 2,
          "service_name": "Washing + Ironing",
          "status": "PROCESSING_IN_PROGRESS",
          "item_count": 5,
          "total_amount": 135.00,
          "current_step": "Iron Items",
          "current_step_order": 3,
          "expected_delivery_date": "2025-11-23"
        },
        {
          "service_id": 3,
          "service_name": "Dry Cleaning",
          "status": "PROCESSING_IN_PROGRESS",
          "item_count": 1,
          "total_amount": 150.00,
          "current_step": "Dry Clean",
          "current_step_order": 4,
          "expected_delivery_date": "2025-11-26"
        }
      ],
      "pricing": {
        "estimated_item_count": 11,
        "actual_item_count": 11,
        "estimated_price": 335.00,
        "final_price": 335.00,
        "discount_amount": 0.00
      },
      "timeline": [
        {
          "status": "BOOKING_CREATED",
          "timestamp": "2025-11-20T10:00:00Z",
          "message": "Order created"
        },
        {
          "status": "PICKUP_SCHEDULED",
          "timestamp": "2025-11-20T10:00:00Z",
          "message": "Pickup scheduled for 2025-11-21 10:30 AM"
        },
        {
          "status": "PICKED_UP",
          "timestamp": "2025-11-21T10:35:00Z",
          "message": "Items picked up. Count verified: 11 items"
        },
        {
          "status": "PROCESSING_IN_PROGRESS",
          "timestamp": "2025-11-21T11:00:00Z",
          "message": "Processing started"
        }
      ],
      "created_at": "2025-11-20T10:00:00Z",
      "updated_at": "2025-11-21T11:00:00Z"
    }
  }
}
```

---

### 6.3 Update Order Count (Vendor)

**Endpoint:** `POST /api/v1/orders/{order_id}/update-count`

**Description:** Vendor updates actual item count at pickup

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body:**
```json
{
  "actual_items": [
    {
      "service_id": 1,
      "rate_card_item_id": 1,
      "item_name": "Shirt",
      "quantity": 7
    },
    {
      "service_id": 2,
      "rate_card_item_id": 4,
      "item_name": "Shirt",
      "quantity": 3
    }
  ],
  "vendor_notes": "Found 2 additional shirts"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Count updated. Awaiting resident approval",
  "data": {
    "order_id": "uuid-v4",
    "status": "COUNT_APPROVAL_PENDING",
    "count_comparison": {
      "original": {
        "total_items": 11,
        "total_amount": 335.00
      },
      "updated": {
        "total_items": 13,
        "total_amount": 355.00
      },
      "difference": {
        "items": 2,
        "amount": 20.00
      }
    },
    "approval_deadline": "2025-11-21T12:35:00Z",
    "auto_approve_in": "2 hours",
    "notification_sent": true
  }
}
```

---

### 6.4 Approve/Reject Count Update (Resident)

**Endpoint:** `POST /api/v1/orders/{order_id}/approve-count`

**Description:** Resident approves or questions the updated count

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body (Approve):**
```json
{
  "action": "APPROVE"
}
```

**Request Body (Question):**
```json
{
  "action": "QUESTION",
  "message": "I don't think there are that many shirts"
}
```

**Response (200 OK - Approved):**
```json
{
  "success": true,
  "message": "Count approved. Order proceeding",
  "data": {
    "order_id": "uuid-v4",
    "status": "PICKED_UP",
    "final_price": 355.00,
    "actual_item_count": 13,
    "count_approved_at": "2025-11-21T10:40:00Z"
  }
}
```

---

## 7. Workflow Management APIs

### 7.1 Get Service Workflow

**Endpoint:** `GET /api/v1/services/{service_id}/workflow`

**Description:** Get workflow steps for a specific service type

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "service": {
      "service_id": 1,
      "service_name": "Ironing Only",
      "service_key": "IRONING"
    },
    "workflow": {
      "template_id": 1,
      "template_name": "Standard Ironing Workflow",
      "is_default": true,
      "steps": [
        {
          "step_id": 1,
          "step_name": "Pickup Items",
          "step_key": "pickup",
          "step_order": 1,
          "is_required": true,
          "requires_photo": false,
          "estimated_duration_minutes": 15,
          "order_status_on_complete": "PICKUP_IN_PROGRESS"
        },
        {
          "step_id": 2,
          "step_name": "Count Items",
          "step_key": "count",
          "step_order": 2,
          "is_required": true,
          "requires_photo": true,
          "estimated_duration_minutes": 10,
          "order_status_on_complete": "COUNT_APPROVAL_PENDING"
        },
        {
          "step_id": 3,
          "step_name": "Iron Items",
          "step_key": "iron",
          "step_order": 3,
          "is_required": true,
          "requires_photo": false,
          "estimated_duration_minutes": 60,
          "order_status_on_complete": "PROCESSING_IN_PROGRESS"
        },
        {
          "step_id": 4,
          "step_name": "Quality Check",
          "step_key": "quality_check",
          "step_order": 4,
          "is_required": true,
          "requires_photo": false,
          "estimated_duration_minutes": 10,
          "order_status_on_complete": "READY_FOR_DELIVERY"
        },
        {
          "step_id": 5,
          "step_name": "Deliver Items",
          "step_key": "deliver",
          "step_order": 5,
          "is_required": true,
          "requires_photo": true,
          "estimated_duration_minutes": 15,
          "order_status_on_complete": "DELIVERED"
        }
      ]
    }
  }
}
```

---

### 7.2 Get Order Workflow Progress

**Endpoint:** `GET /api/v1/orders/{order_id}/workflow-progress`

**Description:** Get workflow progress for all services in an order

**Headers:**
```
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "order_id": "uuid-v4",
    "order_number": "ORD20251120000001",
    "status": "PROCESSING_IN_PROGRESS",
    "services": [
      {
        "service_id": 1,
        "service_name": "Ironing Only",
        "item_count": 5,
        "current_step_order": 4,
        "total_steps": 5,
        "progress_percentage": 80,
        "workflow_steps": [
          {
            "step_id": 1,
            "step_name": "Pickup Items",
            "step_order": 1,
            "status": "COMPLETED",
            "started_at": "2025-11-21T10:30:00Z",
            "completed_at": "2025-11-21T10:35:00Z",
            "duration_minutes": 5,
            "completed_by": "uuid-vendor-id",
            "photos": [],
            "notes": null
          },
          {
            "step_id": 2,
            "step_name": "Count Items",
            "step_order": 2,
            "status": "COMPLETED",
            "started_at": "2025-11-21T10:35:00Z",
            "completed_at": "2025-11-21T10:40:00Z",
            "duration_minutes": 5,
            "completed_by": "uuid-vendor-id",
            "photos": ["https://..."],
            "notes": "Count verified: 5 items"
          },
          {
            "step_id": 3,
            "step_name": "Iron Items",
            "step_order": 3,
            "status": "COMPLETED",
            "started_at": "2025-11-21T11:00:00Z",
            "completed_at": "2025-11-21T17:00:00Z",
            "duration_minutes": 360,
            "completed_by": "uuid-vendor-id"
          },
          {
            "step_id": 4,
            "step_name": "Quality Check",
            "step_order": 4,
            "status": "IN_PROGRESS",
            "started_at": "2025-11-21T17:00:00Z",
            "completed_at": null,
            "estimated_completion": "2025-11-21T17:10:00Z"
          },
          {
            "step_id": 5,
            "step_name": "Deliver Items",
            "step_order": 5,
            "status": "PENDING",
            "started_at": null,
            "completed_at": null
          }
        ]
      },
      {
        "service_id": 2,
        "service_name": "Washing + Ironing",
        "item_count": 5,
        "current_step_order": 3,
        "total_steps": 5,
        "progress_percentage": 60,
        "workflow_steps": [
          {
            "step_id": 6,
            "step_name": "Pickup Items",
            "step_order": 1,
            "status": "COMPLETED",
            "completed_at": "2025-11-21T10:35:00Z"
          },
          {
            "step_id": 7,
            "step_name": "Count Items",
            "step_order": 2,
            "status": "COMPLETED",
            "completed_at": "2025-11-21T10:40:00Z"
          },
          {
            "step_id": 8,
            "step_name": "Wash Items",
            "step_order": 3,
            "status": "COMPLETED",
            "completed_at": "2025-11-21T14:00:00Z"
          },
          {
            "step_id": 9,
            "step_name": "Iron Items",
            "step_order": 4,
            "status": "IN_PROGRESS",
            "started_at": "2025-11-21T14:00:00Z"
          },
          {
            "step_id": 10,
            "step_name": "Quality Check",
            "step_order": 5,
            "status": "PENDING"
          }
        ]
      }
    ]
  }
}
```

---

### 7.3 Complete Workflow Step (Vendor)

**Endpoint:** `POST /api/v1/orders/{order_id}/workflow/{service_id}/complete-step`

**Description:** Mark a workflow step as complete for a specific service

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body:**
```json
{
  "step_id": 3,
  "photos": ["https://photo1.jpg", "https://photo2.jpg"],
  "notes": "All items ironed successfully",
  "signature_url": null
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Workflow step completed",
  "data": {
    "order_id": "uuid-v4",
    "service_id": 1,
    "completed_step": {
      "step_id": 3,
      "step_name": "Iron Items",
      "completed_at": "2025-11-21T17:00:00Z",
      "duration_minutes": 360
    },
    "next_step": {
      "step_id": 4,
      "step_name": "Quality Check",
      "step_order": 4,
      "status": "IN_PROGRESS",
      "started_at": "2025-11-21T17:00:00Z"
    },
    "service_status": "PROCESSING_IN_PROGRESS",
    "order_status": "PROCESSING_IN_PROGRESS",
    "is_final_step": false,
    "notification_sent": true
  }
}
```

---

## 8. Payment APIs

### 8.1 Get Payment Details

**Endpoint:** `GET /api/v1/orders/{order_id}/payment`

**Description:** Get payment information for an order

**Headers:**
```
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "order_id": "uuid-v4",
    "order_number": "ORD20251120000001",
    "payment_status": "PENDING",
    "amount_details": {
      "estimated_price": 335.00,
      "final_price": 355.00,
      "discount_amount": 0.00,
      "amount_due": 355.00
    },
    "payment_methods": [
      {
        "method": "UPI",
        "enabled": true,
        "details": {
          "vendor_upi_id": "perfectpress@paytm",
          "vendor_name": "Perfect Press"
        }
      },
      {
        "method": "CASH",
        "enabled": true,
        "details": {
          "note": "Pay directly to vendor on delivery"
        }
      }
    ],
    "payment_deadline": null,
    "order_status": "DELIVERED",
    "delivered_at": "2025-11-26T15:00:00Z"
  }
}
```

---

### 8.2 Record Payment (Resident)

**Endpoint:** `POST /api/v1/orders/{order_id}/payment`

**Description:** Record payment made by resident

**Headers:**
```
Authorization: Bearer {access_token}
```

**Request Body (UPI):**
```json
{
  "payment_method": "UPI",
  "amount": 355.00,
  "upi_transaction_id": "UPI123456789",
  "upi_vpa": "resident@paytm",
  "payment_screenshot_url": "https://..."
}
```

**Request Body (Cash):**
```json
{
  "payment_method": "CASH",
  "amount": 355.00,
  "notes": "Paid in cash on delivery"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Payment recorded successfully",
  "data": {
    "payment_id": "uuid-v4",
    "order_id": "uuid-v4",
    "amount": 355.00,
    "payment_method": "UPI",
    "status": "COMPLETED",
    "upi_transaction_id": "UPI123456789",
    "paid_at": "2025-11-26T15:05:00Z",
    "order_status": "COMPLETED",
    "next_steps": {
      "message": "Please rate your experience",
      "can_rate": true
    }
  }
}
```

---

## 9. User Profile Management APIs

### 9.1 Get User Profile

**Endpoint:** `GET /api/v1/users/{user_id}/profile`

**Description:** Retrieve complete user profile information

**Headers:**
```
Authorization: Bearer {access_token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "user_id": "uuid-v4",
    "full_name": "Ramesh Kumar",
    "phone": "+919876543210",
    "phone_verified": true,
    "email": "ramesh@example.com",
    "email_verified": true,
    "profile_photo_url": "https://...",
    "user_type": "RESIDENT",
    "created_at": "2025-01-15T10:30:00Z",
    "updated_at": "2025-11-20T14:25:00Z",
    "active_society": {
      "society_id": 1,
      "society_name": "Maple Gardens",
      "unit_type": "FLAT",
      "flat_number": "A-404"
    }
  }
}
```

**Response (200 OK - No Email or Phone):**
```json
{
  "success": true,
  "data": {
    "user_id": "uuid-v4",
    "full_name": "Priya Sharma",
    "phone": null,
    "phone_verified": false,
    "email": null,
    "email_verified": false,
    "profile_photo_url": "https://...",
    "user_type": "RESIDENT",
    "created_at": "2025-11-20T10:30:00Z",
    "updated_at": "2025-11-20T10:30:00Z",
    "needs_contact_update": true,
    "message": "Please add email or phone for better account security"
  }
}
```

---

### 9.2 Update Profile (Basic Info)

**Endpoint:** `PUT /api/v1/users/{user_id}/profile`

**Description:** Update basic profile information (name, photo). Does not require verification.

**Headers:**
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "full_name": "Ramesh Kumar Sharma",
  "profile_photo_url": "https://cdn.example.com/photos/new-photo.jpg"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Profile updated successfully",
  "data": {
    "user_id": "uuid-v4",
    "full_name": "Ramesh Kumar Sharma",
    "profile_photo_url": "https://cdn.example.com/photos/new-photo.jpg",
    "updated_at": "2025-11-20T15:00:00Z"
  }
}
```

---

### 9.3 Request Email Update

**Endpoint:** `POST /api/v1/users/{user_id}/update-email`

**Description:** Initiate email add/update process. Sends verification OTP to the new email address.

**Headers:**
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "new_email": "ramesh.new@example.com"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Verification OTP sent to ramesh.new@example.com",
  "data": {
    "verification_id": "uuid-v4",
    "email": "ramesh.new@example.com",
    "otp_expires_at": "2025-11-20T15:10:00Z",
    "masked_email": "ra****@example.com",
    "next_step": "verify_email_otp"
  }
}
```

**Response (400 Bad Request - Email Already Exists):**
```json
{
  "success": false,
  "error": {
    "code": "EMAIL_ALREADY_EXISTS",
    "message": "This email is already registered with another account"
  }
}
```

---

### 9.4 Verify Email Update

**Endpoint:** `POST /api/v1/users/{user_id}/verify-email`

**Description:** Verify and complete email update using OTP sent to new email

**Headers:**
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "verification_id": "uuid-v4",
  "otp": "123456"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Email verified and updated successfully",
  "data": {
    "user_id": "uuid-v4",
    "email": "ramesh.new@example.com",
    "email_verified": true,
    "updated_at": "2025-11-20T15:08:00Z"
  }
}
```

**Response (400 Bad Request - Invalid OTP):**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_OTP",
    "message": "The OTP entered is incorrect",
    "details": {
      "attempts_remaining": 2,
      "can_resend": true
    }
  }
}
```

**Response (400 Bad Request - Expired OTP):**
```json
{
  "success": false,
  "error": {
    "code": "OTP_EXPIRED",
    "message": "OTP has expired. Please request a new one",
    "details": {
      "can_resend": true
    }
  }
}
```

---

### 9.5 Request Phone Update

**Endpoint:** `POST /api/v1/users/{user_id}/update-phone`

**Description:** Initiate phone add/update process. Sends verification OTP to the new phone number.

**Headers:**
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "new_phone": "+919876543299"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Verification OTP sent to +919876543299",
  "data": {
    "verification_id": "uuid-v4",
    "phone": "+919876543299",
    "otp_expires_at": "2025-11-20T15:10:00Z",
    "masked_phone": "+91****43299",
    "next_step": "verify_phone_otp"
  }
}
```

**Response (400 Bad Request - Phone Already Exists):**
```json
{
  "success": false,
  "error": {
    "code": "PHONE_ALREADY_EXISTS",
    "message": "This phone number is already registered with another account"
  }
}
```

---

### 9.6 Verify Phone Update

**Endpoint:** `POST /api/v1/users/{user_id}/verify-phone`

**Description:** Verify and complete phone update using OTP sent to new phone number

**Headers:**
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "verification_id": "uuid-v4",
  "otp": "123456"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Phone number verified and updated successfully",
  "data": {
    "user_id": "uuid-v4",
    "phone": "+919876543299",
    "phone_verified": true,
    "updated_at": "2025-11-20T15:08:00Z"
  }
}
```

**Response (400 Bad Request - Invalid OTP):**
```json
{
  "success": false,
  "error": {
    "code": "INVALID_OTP",
    "message": "The OTP entered is incorrect",
    "details": {
      "attempts_remaining": 2,
      "can_resend": true
    }
  }
}
```

---

### 9.7 Resend Verification OTP

**Endpoint:** `POST /api/v1/users/{user_id}/resend-verification-otp`

**Description:** Resend OTP for email or phone verification

**Headers:**
```
Authorization: Bearer {access_token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "verification_id": "uuid-v4",
  "type": "EMAIL"
}
```

**Request Body (Phone):**
```json
{
  "verification_id": "uuid-v4",
  "type": "PHONE"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "OTP resent successfully",
  "data": {
    "verification_id": "uuid-v4",
    "type": "EMAIL",
    "masked_contact": "ra****@example.com",
    "otp_expires_at": "2025-11-20T15:15:00Z",
    "can_resend_after": "2025-11-20T15:06:00Z"
  }
}
```

**Response (429 Too Many Requests):**
```json
{
  "success": false,
  "error": {
    "code": "TOO_MANY_REQUESTS",
    "message": "Please wait before requesting another OTP",
    "details": {
      "retry_after_seconds": 45,
      "can_resend_after": "2025-11-20T15:06:00Z"
    }
  }
}
```

---

### 9.8 Remove Contact Information

**Endpoint:** `DELETE /api/v1/users/{user_id}/contact/{type}`

**Description:** Remove email or phone from profile. At least one contact method (email or phone) must remain.

**Headers:**
```
Authorization: Bearer {access_token}
```

**Path Parameters:**
- `type` - Contact type to remove: `email` or `phone`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Email removed successfully",
  "data": {
    "user_id": "uuid-v4",
    "email": null,
    "email_verified": false,
    "phone": "+919876543210",
    "phone_verified": true,
    "updated_at": "2025-11-20T15:20:00Z"
  }
}
```

**Response (400 Bad Request - Last Contact Method):**
```json
{
  "success": false,
  "error": {
    "code": "LAST_CONTACT_METHOD",
    "message": "Cannot remove last contact method. Please add another contact before removing this one",
    "details": {
      "current_contacts": {
        "email": "ramesh@example.com",
        "phone": null
      }
    }
  }
}
```

---

## API Response Codes

### Success Codes
- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `204 No Content` - Successful with no response body

### Client Error Codes
- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Authentication required or failed
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict (e.g., duplicate entry)
- `422 Unprocessable Entity` - Validation failed
- `429 Too Many Requests` - Rate limit exceeded

### Server Error Codes
- `500 Internal Server Error` - Server error
- `503 Service Unavailable` - Service temporarily unavailable

---

## Error Response Format

All error responses follow this structure:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {
      "field": "specific_field",
      "reason": "validation failure reason"
    }
  }
}
```

---

## Common Headers

### Request Headers
```
Authorization: Bearer {access_token}
Content-Type: application/json
Accept: application/json
X-API-Version: v1
```

### Response Headers
```
Content-Type: application/json
X-Request-ID: uuid-v4
X-Rate-Limit-Remaining: 100
X-Rate-Limit-Reset: 1700000000
```

---

**End of API Specification Document**
