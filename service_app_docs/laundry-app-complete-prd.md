# Product Requirements Document (PRD): Multi-Society Laundry Services App
## Complete Document

**Version:** 2.1
**Last Updated:** November 17, 2025
**Document Owner:** Product Team

---

## Table of Contents

### Part 1: Resident & Laundry Person Flows
1. [Executive Summary](#1-executive-summary)
2. [Order Lifecycle](#2-order-lifecycle)
3. [Resident Flows](#3-resident-flows)
4. [Laundry Person Flows](#4-laundry-person-flows)

### Part 2: Admin Flows, Payment Processing & Dispute Resolution
5. [Admin Flows](#5-admin-flows)
6. [Payment Processing System](#6-payment-processing-system)
7. [Dispute Resolution System](#7-dispute-resolution-system)
8. [Trust & Safety Framework](#8-trust--safety-framework)

### Part 3: Payment & Settlement Model
9. [Payment Flow Architecture](#9-payment-flow-architecture)
10. [Payment Collection Methods](#10-payment-collection-methods)
11. [Settlement Model](#11-settlement-model)
12. [Platform Revenue Model](#12-platform-revenue-model)
13. [Payment Enforcement & Collection](#13-payment-enforcement--collection)
14. [Financial Reconciliation](#14-financial-reconciliation)
15. [Edge Cases & Scenarios](#15-edge-cases--scenarios)
16. [Summary & Recommendations](#16-summary--recommendations)

---

# Part 1: Resident & Laundry Person Flows

## 1. Executive Summary

### 1.1 Objective
Deliver a user-friendly laundry services management platform that connects residents with laundry service providers in their residential societies, enabling transparent pricing, real-time tracking, and reliable service delivery across multiple service types (ironing, washing, dry cleaning).

### 1.2 Key Principles

**Simplicity:**
- Maximum 2 taps for any core action
- Plain language, no jargon
- Smart defaults that learn from usage

**Transparency:**
- Clear pricing before booking
- Real-time order status
- Count verification with approval

**Trust:**
- Rating system for quality
- Photo proof options
- Fair dispute resolution

**Leniency:**
- No-penalty cancellations
- Flexible policies
- Trust-based approach

### 1.3 Scope

This document covers complete user flows for:
- Residents (booking, tracking, payment, ratings, issues)
- Laundry Persons (registration, task management, settlements) (Future scope: Other services)
- Admin flows
- Payment processing details
- Dispute resolution system
- Settlement mechanisms

---

## 2. Order Lifecycle

### 2.1 Complete Order States

**1. BOOKING CREATED**
- Resident completes booking with item count and pickup time
- Estimated price calculated based on rate card
- Order ID generated

**2. PICKUP SCHEDULED**
- Laundry person notified of new order
- Both parties receive confirmation with pickup time
- 30-minute reminder scheduled before pickup

**3. PICKUP IN PROGRESS**
- Laundry person marks start of pickup
- En route to resident's address

**4. COUNT APPROVAL PENDING** (conditional)
- Laundry person updates actual count at pickup if different from booking
- Resident receives notification with price adjustment
- 2-hour approval window
- Auto-approves after timeout

**5. PICKED UP**
- Items collected and count verified/approved
- Final price locked
- Ready for service to begin

**6. PROCESSING IN PROGRESS**
- Laundry person marks service started (washing/ironing/dry cleaning)
- Work actively happening

**7. READY FOR DELIVERY**
- Service completed
- Ready to be delivered back

**8. OUT FOR DELIVERY**
- Laundry person marks delivery started
- En route to resident's address

**9. DELIVERED**
- Items returned to resident
- Delivery confirmed with timestamp
- Payment flow triggered

**10. COMPLETED**
- Payment received and confirmed
- Rating submitted by resident
- Order archived

**ALTERNATE STATES:**
- **CANCELLED**: Order cancelled before or shortly after pickup
- **DISPUTED**: Issue reported by resident
- **ON HOLD**: Temporarily paused by resident request

### 2.2 State Transition Rules

**Cannot Cancel After:**
- PICKED UP state (items in laundry person's possession)
- Must complete order or coordinate return

**Cannot Skip States:**
- Must go through IRONING IN PROGRESS
- Cannot jump from PICKED UP to DELIVERED

**Count Approval:**
- Only triggered if actual count differs from booking
- Otherwise transitions directly from PICKUP IN PROGRESS to PICKED UP

**Timing Requirements:**
- Pickup: Within 48 hours of booking
- Count approval: 2-hour window
- Ironing: Typically 24-48 hours
- Delivery: Within 48 hours of ready state

---

## 3. Resident Flows

### 3.1 Registration & Onboarding

**Flow:**

1. User enters phone number
2. Receives OTP via SMS (< 5 second delivery)
3. Enters OTP to verify phone
4. System checks society roster database
   - **If phone found in roster:** Auto-assigns society and flat, asks for name confirmation
   - **If phone not found:** User selects society from dropdown, enters flat number, submits for society admin verification
5. Profile created

**Verification Pending State:**
- User can browse service providers and rate cards
- Cannot place orders until verified by admin
- Typically verified within 24 hours

**Edge Cases:**
- Multiple properties: User can add additional societies later via profile
- Moving flats: Update triggers re-verification
- Roster mismatch: Admin manually corrects and approves

---

### 3.2 Browse Service Categories

**Purpose:** Allow residents to discover all available service types

**Flow:**

1. **Home Screen** displays service category cards:
   ```
   ┌─────────────────────────────┐
   │ 👔 Laundry Services         │
   │ Iron, Wash, Dry Clean       │
   │ 50 providers near you       │
   └─────────────────────────────┘

   ┌─────────────────────────────┐
   │ 🚗 Vehicle Services         │
   │ Coming Soon                 │
   └─────────────────────────────┘

   ┌─────────────────────────────┐
   │ 🏡 Home Services            │
   │ Coming Soon                 │
   └─────────────────────────────┘
   ```

2. User taps a category (currently only Laundry is active)

3. **Category Details Screen** shows:
   - All subcategories (Ironing, Washing, Dry Cleaning, etc.)
   - Number of providers offering each subcategory
   - Typical price range per subcategory
   - Average turnaround time

4. User selects a specific subcategory or views all providers in category

**Future State (Multi-Category):**
- When vehicle/home services launch, same flow applies
- Users can browse across categories
- Favorites saved per category
- Order history grouped by category

---

### 3.3 Browse Service Providers (Within Category)

**Flow:**

1. User lands on home screen showing service providers serving their society

2. **Filter by Service Type** (top of screen):
   - All Services (default)
   - Ironing Only
   - Washing + Ironing
   - Dry Cleaning
   - Washing Only

3. List displays each service provider with:
   - Name, rating, order count
   - **Services offered badges**: [Ironing] [Washing] [Dry Cleaning]
   - **Starting prices**: "From ₹10/piece (Ironing)" or "From ₹25/piece (Wash+Iron)"
   - Estimated delivery time
   - Availability status

4. User can tap any provider to view details:
   - **Complete rate card** (tabbed by service type):
     - Ironing tab
     - Washing + Ironing tab
     - Dry Cleaning tab
     - Washing Only tab
   - Recent reviews (filtered by service type optionally)
   - Service description
   - Turnaround time per service
   - Contact options

5. User selects service provider to proceed with booking

**Sorting:**
- Primary: Rating (highest first)
- Secondary: Delivery time (fastest first)
- Tertiary: Total orders (reliability indicator)

**Service Badges:**
- Each provider shows colored badges for services offered
- Makes it easy to see at a glance what services are available

---

### 3.3 Create New Order

**Flow:**

**Step 1: Add Items with Service Selection**

User builds order by adding items with specific services:

**Option A: Quick Add (Faster)**
- Select service type first (Ironing, Washing+Ironing, Dry Cleaning, Washing Only)
- Enter total item count for that service
- Example: "10 items for Washing + Ironing"
- Can add multiple service types:
  - 10 items for Washing + Ironing
  - 5 items for Ironing Only
  - 2 items for Dry Cleaning
- System calculates estimate based on average price per service type

**Option B: Detailed Item Selection (More Accurate)**
- Browse complete rate card showing all services
- Rate card organized by service type tabs or sections:

  **🔵 Ironing Only**
  - Shirt ₹10 → [Add] or [+/-]
  - Pants ₹15 → [Add]
  - Saree ₹30 → [Add]

  **🟢 Washing + Ironing**
  - Shirt ₹25 → [Add]
  - Pants ₹30 → [Add]
  - Saree ₹60 → [Add]

  **🟡 Dry Cleaning**
  - Shirt ₹80 → [Add]
  - Blazer ₹150 → [Add]
  - Saree ₹200 → [Add]

- User selects items from any service type
- Each item selection includes service type + quantity
- Running total updates in real-time

**Example Mixed Order:**
```
Shopping Cart:
├─ 5 Shirts (Washing + Ironing) @ ₹25 = ₹125
├─ 3 Pants (Washing + Ironing) @ ₹30 = ₹90
├─ 2 Shirts (Ironing Only) @ ₹10 = ₹20
├─ 1 Blazer (Dry Cleaning) @ ₹150 = ₹150
└─ 2 Sarees (Dry Cleaning) @ ₹200 = ₹400

Total: ₹785
Service Types: 3 (Washing+Iron, Iron, Dry Clean)
Total Items: 13
```

**Step 2: Review Cart**
- Shows items grouped by service type
- Breakdown of costs per service
- Total estimated price
- Item count per service
- Can edit quantities or remove items
- Disclaimer: "Final amount confirmed after pickup count verification"

**Step 3: Select Pickup Time**

System suggests optimal pickup time based on:
- Service provider's existing schedule
- Daily capacity limits
- Multiple service types in order (uses longest turnaround as base)
- Time of booking (if before 10 AM, suggests same day; otherwise next day)
- Provides 3 alternative time slots
- User can select custom time if none suitable

**Step 4: Expected Delivery Date**

System calculates delivery based on longest turnaround service:
- If order has Dry Cleaning (3-5 days): Delivery in 5 days
- If order has Washing + Ironing (2-3 days) but no Dry Cleaning: Delivery in 3 days
- If order has only Ironing/Washing: Delivery in 1-2 days

**Option for Partial Delivery:**
- User can request partial deliveries (e.g., get ironing items first, dry cleaning later)
- Creates separate delivery tracking per service type
- Additional coordination required

**Step 5: Confirm Pickup Address**
- Auto-filled with user's flat and society
- Option to add special instructions (floor number, gate code, etc.)
- Option to add alternate contact number
- Note if multiple service types: "All items will be picked up together"

**Step 6: Review and Confirm**
- Summary displays:
  - Service provider name and rating
  - **Items by service type** (expandable):
    - Washing + Ironing: 8 items, ₹215
    - Ironing Only: 2 items, ₹20
    - Dry Cleaning: 3 items, ₹550
  - **Total**: 13 items, ₹785
  - Pickup date/time and address
  - **Expected delivery date**: Based on slowest service (e.g., 5 days for dry cleaning)
  - Breakdown showing estimated delivery per service type
- Disclaimer: "Final amount confirmed after pickup count verification"
- User confirms to create order

**Post-Confirmation:**
- **Single order** created with multiple service types
- Order status PICKUP SCHEDULED
- Service provider notified immediately
- Both parties receive confirmation
- 30-minute reminder scheduled
- Order summary shows all service types included

**Business Logic:**
- Quick count estimate = Sum of (Items × Average price per service type)
- Detailed count = Sum of (Quantity × Rate for each item + service combination)
- Pickup time must be minimum 2 hours in future
- System prevents double-booking same time slot
- **Delivery time = Maximum turnaround among all services in order**
  - If order has Dry Cleaning: 3-5 days
  - Else if order has Washing + Ironing: 2-3 days
  - Else if order has Washing Only: 1-2 days
  - Else (Ironing only): 1-2 days
- Partial delivery option extends timeline but provides better customer experience

---

### 3.4 Count Update Approval

**Trigger:**
Laundry person arrives at pickup and actual item count differs from what resident booked

**Flow:**

1. Laundry person enters actual count in their app
2. System calculates new price based on rate card
3. Resident receives immediate notification:
   - Shows original booking (10 items, ₹200)
   - Shows actual count (12 items, ₹240)
   - Shows difference (+2 items, +₹40)
4. Resident has three options:

**Option A: Approve**
- One-tap approval
- Final price locked at new amount
- Order status → PICKED UP
- Ironing can begin

**Option B: Question**
- Opens communication with laundry person (call or chat)
- Laundry person can explain or re-count
- If adjustment needed, updated count sent for re-approval
- Maximum 2 adjustment cycles
- After 2 cycles, option to escalate to admin

**Option C: Wait (No Action)**
- After 2 hours with no response: Auto-approve
- Notification sent: "Count auto-approved at ₹240"
- Order proceeds normally
- Resident can dispute within 24 hours of delivery if needed

**Edge Cases:**
- Resident unavailable: Approval request queued, auto-approves after 2 hours
- Major discrepancy (>50% difference): Flagged for admin review
- Multiple updates: Second update within 1 hour requires admin approval

---

### 3.5 Track Order

**Flow:**

1. User accesses order from home screen or orders list
2. Tracking view displays:
   - Current status with visual progress indicator
   - Complete timeline of all states with timestamps
   - Expected completion time for current state
   - Next steps in the process
3. Real-time updates via push notifications for each state change
4. User can take actions based on current state:
   - Call or message laundry person
   - Cancel order (if eligible)
   - View order details

**Status Visibility:**
- Each completed state shows timestamp and any relevant notes
- Current state shows start time and expected completion
- Future states show estimated times
- Delays or issues flagged with alerts

**Notifications Sent:**
- PICKUP SCHEDULED: "Pickup confirmed for tomorrow 10:30 AM"
- PICKED UP: "Your clothes picked up! Ironing starts soon"
- IRONING IN PROGRESS: "Ironing in progress. Ready by tomorrow afternoon"
- READY FOR DELIVERY: "Your order is ready for delivery"
- DELIVERED: "Order delivered! Pay now to complete"

---

### 3.6 Receive Delivery

**Flow:**

**Scenario 1: In-Person Delivery**
1. Laundry person arrives and hands over items to resident
2. Resident can visually verify items
3. Laundry person marks order as delivered
4. Status changes to DELIVERED
5. Payment flow triggered immediately

**Scenario 2: Delivery to Guard/Neighbor**
1. Resident not available at time of delivery
2. Laundry person leaves items with security guard or neighbor
3. Takes photo proof of handover (Phase 2)
4. Marks delivered with note: "Left with security guard at main gate"
5. Resident receives notification with delivery details and photo
6. Resident options:
   - Confirm receipt after collecting
   - Report "Not received yet" if actually not delivered
7. Payment flow triggered

**Item Verification:**
- Resident encouraged to verify count immediately
- Any issues should be reported before payment
- After payment, disputes require more evidence

---

### 3.7 Payment

**Flow:**

After delivery, resident sees payment options:

**UPI Payment:**
1. User selects "Pay via UPI"
2. System generates payment link with amount and order details
3. Opens user's preferred payment app (Google Pay, PhonePe, etc.)
4. User completes payment in payment app
5. System waits for payment confirmation (30-second timeout)
6. On success:
   - Payment marked complete
   - Both parties notified
   - Digital receipt generated
   - Rating prompt appears
7. On failure:
   - User can retry (max 3 attempts)
   - After 3 failures, suggests cash payment option

**Cash Payment:**
1. User selects "Mark as Cash Paid"
2. Confirms they paid cash to laundry person
3. Laundry person receives verification request
4. Order shows "Payment Pending Confirmation" until laundry person confirms
5. Once confirmed, status changes to COMPLETED

**Pay Later:**
1. User can defer payment
2. After 48 hours: Soft reminder notification
3. Restrictions after 48 hours:
   - Can still book with SAME laundry person (trust-based)
   - Cannot book with NEW laundry persons
4. After 2+ unpaid orders: Account suspended until dues cleared

**Offline Payment:**
- If user offline during payment: Request queued
- Auto-retries when connection restored
- Notification sent when completed

---

### 3.8 Rate & Review

**Flow:**

1. After payment completed, rating prompt appears
2. User rates experience 1-5 stars (tap to select)
3. Optional comment field provided
4. If 1-2 stars: Comment becomes mandatory
5. User submits rating
6. Rating added to laundry person's profile

**Rules:**
- Can only rate after delivery + payment
- One rating per order
- Cannot change rating after submission
- Rating affects laundry person's average and visibility

**Abuse Prevention:**
- System tracks if resident gives 5+ consecutive 1-star ratings
- Flagged for admin review
- Patterns of abuse may result in ratings being discounted

---

### 3.9 Report Issue

**Flow:**

1. User accesses order and selects "Report Issue" (available after delivery, within 7 days)
2. Selects issue type:
   - Item(s) missing
   - Item damaged/torn
   - Stain not removed
   - Poor ironing quality
   - Wrong items returned
   - Other
3. Uploads up to 3 photos as evidence (optional but encouraged)
4. Describes issue in text (500 character limit)
5. Selects expected resolution:
   - Full refund
   - Partial refund (specify amount)
   - Re-iron items
   - Replace item
   - Just informing (no action)
6. Submits issue report

**After Submission:**
- Ticket created with unique ID
- Laundry person notified immediately
- 24-hour window for laundry person to respond
- Issue tracking view shows real-time status
- Resident receives updates as laundry person responds

**Resolution Paths:**
- Laundry person accepts and offers compensation
- Laundry person promises to return missing item
- Laundry person disputes with counter-evidence → Admin escalation
- No response after 24 hours → Auto-escalate to admin

---

### 3.10 Cancel Order

**Flow:**

**Before Pickup:**
1. User selects order and taps "Cancel Order"
2. Optional: Select cancellation reason
3. Confirms cancellation
4. Order immediately cancelled
5. Laundry person notified
6. No penalties or restrictions

**After Pickup, Before Ironing:**
1. User selects order and taps "Cancel Order"
2. System shows: "Items already picked up. Need to arrange return."
3. Options presented:
   - "I'll collect from laundry person"
   - "Request return delivery"
4. If return delivery: Laundry person can accept or decline
5. Once return arranged, order cancelled
6. No charges

**After Ironing Started:**
1. "Cancel Order" button disabled
2. Message shows: "Cannot cancel - ironing in progress"
3. User must contact laundry person directly
4. Typical resolution: Complete order, pay for work done

**Laundry Person Cancellation:**
- Before pickup: Can cancel with reason (emergency, overbooked, etc.)
- After pickup: Cannot cancel (must complete or find substitute)
- System tracks cancellation frequency for quality monitoring

**Pattern Tracking (No Penalties):**
- Laundry person: 10+ cancellations/month → Admin reaches out supportively
- Resident: 15+ cancellations/month → Gentle awareness reminder
- No blocking, just pattern awareness

---

### 3.11 Multi-Society Management (Phase 2)

**Flow:**

1. User accesses profile and selects "Manage Locations"
2. Sees current verified societies/flats
3. Taps "Add New Location"
4. Selects society from dropdown
5. Enters flat number
6. Submits for verification
7. Admin verifies (typically < 24 hours)
8. Once verified, appears in location list

**When Ordering:**
1. Home screen shows current society: "Ordering from: Maple Society"
2. User taps society name to switch
3. Modal shows all verified societies
4. User selects different society
5. Laundry person list updates to show only those serving selected society
6. User proceeds with order as normal

---

### 3.12 Repeat Order (Phase 2)

**Flow:**

1. User taps "Repeat Last Order" from home screen
2. System checks last order's laundry person
   - **If available:** Pre-fills order with last order's details (items, laundry person, pickup preferences)
   - **If unavailable:** Suggests similar alternatives with comparable pricing
3. User can edit any details before confirming
4. Rest of flow identical to new order creation
5. Faster booking (saves 2-3 steps)

---

## 4. Service Provider (Vendor) Flows

### 4.1 Registration & Onboarding

**Flow:**

1. User enters phone number
2. Receives and enters OTP for verification
3. Fills personal information form:
   - Full name
   - Business name (optional)
   - Store address (complete address with landmark)
   - Store photo (Phase 2)

4. **Selects Service Categories to Offer:**
   ```
   Which services do you provide?

   👔 Laundry Services
   ☑ Ironing Only
   ☑ Washing + Ironing
   ☐ Dry Cleaning
   ☐ Washing Only

   🚗 Vehicle Services (Coming Soon)
   ☐ Car Wash
   ☐ Bike Wash

   🏡 Home Services (Coming Soon)
   ☐ Gardening
   ☐ Plumbing
   ```

   - Can select multiple services
   - Will need to set up separate rate card for each selected service
   - Can add/remove services later from profile

5. Selects societies to serve:
   - Search and select from dropdown (1-5 societies recommended)
   - Each society requires separate admin approval

6. Identity verification (Phase 2):
   - Upload photo ID (Aadhaar/License/PAN)
   - Selfie at store location

7. Submits application

**Post-Submission:**
- Status: "Pending Approval"
- Can set up rate cards in draft mode for each selected service
- Cannot receive orders until approved
- Notification when admin approves (typically < 24 hours)

**Approval:**
- Admin reviews and approves for each society separately
- May approve some societies and reject others
- Rejection includes reason and option to reapply
- Once approved, must publish rate card(s) before receiving orders

**Future: Multi-Category Providers**
- Can offer services across multiple categories (e.g., Laundry + Vehicle Services)
- Separate rate cards for each category
- Dashboard groups orders by category
- Analytics show performance per category

---

### 4.2 Rate Card Setup (Per Service Category)

**Flow:**

**Step 1: Select Service Types**

Laundry person selects which services they offer:
- ☐ Ironing Only
- ☐ Washing + Ironing
- ☐ Dry Cleaning
- ☐ Washing Only

Can select multiple service types.

**Step 2: Setup Rate Card per Service Type**

For each selected service, create a rate card:

**Option 1: Template-Based (Recommended)**
1. System shows pre-filled template with common items and standard prices:

   **Ironing Only:**
   - Shirt ₹10, Pants ₹15, Saree ₹30, Bedsheet ₹25, etc.

   **Washing + Ironing:**
   - Shirt ₹25, Pants ₹30, Saree ₹60, Bedsheet ₹50, etc.

   **Dry Cleaning:**
   - Shirt ₹80, Pants ₹100, Blazer ₹150, Saree ₹200, etc.

   **Washing Only:**
   - Shirt ₹20, Pants ₹25, Bedsheet ₹40, etc.

2. Laundry person edits prices to match their rates
3. Can add custom items not in template
4. Can hide items they don't service
5. Previews how rate card appears to customers
6. Publishes rate card

**Option 2: Custom Build**
1. Starts with blank rate card
2. For each service type, adds items one by one:
   - Item name
   - Service type (Ironing/Washing+Ironing/Dry Cleaning/Washing)
   - Price per piece
   - Optional description
3. Previews and publishes

**Option 3: PDF Upload**
1. Uploads existing rate card PDF (max 5MB)
2. PDF shown to customers for reference
3. Must still enter digital pricing for automated billing
4. Publishes

**Rate Card Display:**
- Customers see services grouped by type
- Can filter by service type when browsing
- Clearly labeled: "Ironing: ₹10 | Washing+Ironing: ₹25 | Dry Cleaning: ₹80"

**Multi-Society Rate Cards:**
- Choose at setup: Same rates for all societies OR different rates per society
- If different: Create separate rate card for each society
- Can clone and modify base rate card

**Editing Published Rate Card:**
- Can edit anytime from settings
- Can add/remove service types
- Changes apply only to NEW orders
- Existing orders use original pricing
- Customers notified of rate changes

**Validation:**
- Minimum 5 items required per service type
- Price range: ₹5 to ₹500 per item
- At least one common item must be present (shirt/pants)
- At least one service type must be selected

---

### 4.3 Receive & Accept Orders

**Flow:**

1. Resident books order
2. Laundry person receives immediate notification:
   - Order ID and customer details
   - Pickup time and address
   - Item count and estimated value
3. Order automatically added to task list
4. No explicit "accept" action required (auto-accepted)
5. Can cancel before pickup if needed (with reason)

**Notification Details:**
- "New order from A-404, Maple Society"
- "Pickup: Tomorrow 10:30 AM"
- "12 items estimated, ~₹240"
- Deep link opens order details

---

### 4.4 Task Dashboard

**Flow:**

1. Laundry person opens app to view today's tasks
2. Dashboard organized into sections:
   - **Urgent (next 2 hours):** Immediate pickups and deliveries
   - **Today's Summary:** Total pickups, ironing in progress, deliveries pending
   - **Pending Payments:** Orders delivered but unpaid
3. Can filter by:
   - Society
   - Status (pickup/delivery/ironing)
   - Building/floor (for route optimization)
4. Each task shows key info:
   - Customer name and flat
   - Order ID
   - Item count
   - Amount
   - Scheduled time
   - Action button (Start Pickup/Deliver)

**Task Grouping:**
- For 50+ orders: Groups by building and floor
- Shows count per building: "Building A (5 tasks)"
- Tap to expand and see individual tasks
- Helps with route planning

---

### 4.5 Pickup Flow

**Flow:**

1. Laundry person sees pickup task in dashboard
2. Taps "Start Pickup" when ready
3. Navigation can be opened to customer address (optional)
4. Arrives at customer location
5. Collects items and counts them
6. Enters actual count in app
7. **If count matches booking:**
   - Marks "Pickup Complete"
   - Status changes to PICKED UP
   - Can immediately start ironing
8. **If count differs from booking:**
   - Enters actual count
   - System calculates new price
   - Marks "Pickup Complete"
   - Sends approval request to resident
   - Status: COUNT APPROVAL PENDING
   - Waits for approval before ironing (or 2-hour auto-approval)
9. Optional: Take photo of items (Phase 2)

**Edge Cases:**

**Customer Not Home:**
- Options presented:
  - Call customer
  - Reschedule pickup
  - Cancel order
- Max 2 reschedules allowed
- After 2 reschedules: Order auto-cancels (no penalty to laundry person)

**Cannot Reach Customer:**
- Laundry person can cancel with reason "Cannot reach customer"
- No penalty for laundry person
- Resident notified

**Major Count Discrepancy:**
- If actual count >50% different from booking: System flags for review
- Encourages taking photo proof
- Admin visibility for quality monitoring

---

### 4.6 Update Order Status

**Flow:**

**Mark Ironing Started:**
1. After pickup complete and count approved
2. Laundry person finds order in task list
3. Taps "Start Ironing"
4. Status changes to IRONING IN PROGRESS
5. Resident notified
6. Expected completion time calculated (typically +24 hours)

**Mark Ready for Delivery:**
1. When ironing complete
2. Laundry person taps "Mark Ready"
3. Status changes to READY FOR DELIVERY
4. Order moves to delivery task list
5. Resident notified: "Your order is ready!"
6. Can schedule delivery or wait for customer pickup (in-store)

**Mark Out for Delivery:**
1. Laundry person ready to deliver
2. Taps "Start Delivery"
3. Status changes to OUT FOR DELIVERY
4. Can open navigation to customer address
5. Resident notified: "Your order is on the way"

---

### 4.7 Delivery Flow

**Flow:**

1. Laundry person arrives at customer address
2. Two scenarios:

**Scenario A: Customer Present**
- Hands over items in person
- Customer can verify items
- Laundry person taps "Mark Delivered"
- Selects payment status:
  - "Cash collected" (enter amount)
  - "Will pay online"
- Status changes to DELIVERED
- If cash collected: Order settlement updated immediately

**Scenario B: Customer Not Present**
- Laundry person leaves with guard/neighbor
- Enters delivery details:
  - Delivered to: (Security guard/Neighbor/Family member)
  - Name of person
  - Location (main gate, etc.)
- Takes photo proof (Phase 2)
- Taps "Mark Delivered"
- Resident notified with delivery details and photo
- Resident confirms receipt later

**Partial Delivery (Phase 2):**
- If cannot deliver all items:
  - Enter count of items delivered
  - Add reason for partial delivery
  - Remaining items scheduled for next delivery
  - Price adjusted proportionally

---

### 4.8 Settlement & Money Tracking

**Flow:**

1. Laundry person accesses "Money" section
2. Dashboard shows:
   - This week's earnings (total)
   - Collected amount (UPI + confirmed cash)
   - Pending amount (unpaid orders)
3. Breakdown by payment method:
   - UPI payments (auto-recorded)
   - Cash collected (self-reported)
   - Pending online payments
   - Pending cash payments
4. List of pending payments shows each order:
   - Customer name and flat
   - Order ID and amount
   - Days since delivery
   - Actions: "Mark Cash Received" or "Send Reminder"

**Mark Cash Received:**
1. Laundry person taps "Mark Cash Received" for an order
2. Confirms amount received
3. System updates order status to paid
4. Resident receives confirmation notification
5. If resident disputes: Escalates to admin

**Send Payment Reminder:**
1. Laundry person taps "Send Reminder"
2. Polite notification sent to resident
3. Limited to 1 reminder per 24 hours per order
4. After 7 days unpaid: Option to escalate to admin

**Settlement by Society (Multi-Society):**
- Can filter earnings by society
- See performance per society
- Helps identify profitable areas

**UPI QR Code:**
- Laundry person can generate and share UPI QR code
- Customers can scan to pay directly
- Payments auto-reflected in settlement

---

### 4.9 Respond to Issues

**Flow:**

1. Resident reports issue on an order
2. Laundry person receives immediate notification
3. Views issue details:
   - Issue type (missing item, damage, quality, etc.)
   - Customer description
   - Photos uploaded (if any)
   - Expected resolution
4. Has 24 hours to respond
5. Three response options:

**Option A: Accept & Compensate**
- Acknowledges issue
- Offers refund amount
- Adds explanation
- Submits response
- Resident can accept offer or negotiate
- Once agreed, amount deducted from settlement

**Option B: Item Will Be Returned**
- For missing items
- Sets return/delivery date
- Adds explanation
- Resident waits for return
- If not returned by date: Auto-escalates to admin

**Option C: Dispute**
- Disagrees with claim
- Provides counter-explanation
- Uploads counter-evidence (photos, etc.)
- Automatically escalates to admin for review
- Admin makes final decision

**No Response:**
- If no response within 24 hours: Auto-escalates to admin
- Laundry person receives warning
- Admin reviews with available evidence

---

### 4.10 In-Store Delivery (Phase 2)

**Flow:**

1. Customer walks into store with clothes
2. Laundry person taps "Add In-Store Order"
3. Searches customer:
   - By phone number OR
   - By flat number
4. System finds registered customer
5. Confirms customer details
6. Enters order details:
   - Item count
   - Takes photo of items (optional)
   - Expected ready time
7. Sends approval request to customer
8. Customer receives notification on their phone:
   - Order details and estimated price
   - Options: Approve or Reject
9. **If approved:**
   - Order created with status IRONING IN PROGRESS
   - Flows into normal delivery process
10. **If rejected:**
    - Laundry person notified
    - Can try again with adjusted details or proceed outside app
11. **If no response:**
    - After 4 hours: Auto-reject
    - Laundry person notified
    - Can re-send request or handle outside app

**Trusted Customer Feature (Phase 3):**
- After 5+ successful in-store deliveries
- Laundry person can mark customer as "Trusted"
- Future in-store orders auto-approve up to ₹500
- Customer still notified, can cancel within 1 hour

---

### 4.11 Manage Availability (Phase 2)

**Flow:**

**Set Days Off:**
1. Laundry person accesses "Availability" settings
2. Views calendar with current schedule
3. Taps "Add Day Off"
4. Selects date
5. Chooses reason (weekly off, festival, personal work)
6. System checks for existing orders on that day
7. If orders exist: Prompts to reschedule or complete before day off
8. Confirms day off

**Vacation Mode:**
1. Taps "Vacation Mode"
2. Sets start and end dates
3. During vacation:
   - No new bookings accepted
   - Profile hidden from customer searches
   - Auto-message sent to inquiries
4. Must complete or reschedule pending orders before vacation starts
5. Activates vacation mode

**Daily Capacity Settings:**
1. Sets maximum pickups per day (e.g., 15)
2. Sets maximum deliveries per day (e.g., 20)
3. Blocks specific time slots (lunch break, etc.)
4. System prevents overbooking

---

### 4.12 View Analytics (Phase 3)

**Flow:**

1. Laundry person accesses "Analytics" section
2. Dashboard shows monthly overview:
   - Total orders and earnings
   - Average per order
   - Current rating
   - Trends compared to last month
3. Detailed breakdowns:
   - Orders by day (line graph)
   - Payment method distribution (UPI vs Cash)
   - Top customers (by order frequency)
   - Item type breakdown (most ironed items)
   - Society-wise performance (if serving multiple)
4. Can download monthly report
5. Insights help with:
   - Pricing optimization
   - Capacity planning
   - Customer retention strategies

---

# Part 2: Admin Flows, Payment Processing & Dispute Resolution

## 5. Admin Flows

### 5.1 Admin Roles & Responsibilities

**Society Admin/Manager:**
- Approve/reject laundry person applications
- Manage resident rosters
- Resolve disputes between residents and laundry persons
- Monitor service quality
- Handle escalations

**Capabilities:**
- View all orders in their society
- Communication with all parties
- Financial dispute resolution
- Quality enforcement actions

---

### 5.2 Approve Laundry Person Registration

**Flow:**

1. Laundry person submits registration application
2. Admin receives notification: "New laundry person application"
3. Admin views application details:
   - Personal information (name, phone, business name)
   - Store address and location
   - Societies requested to serve
   - ID proof (Phase 2)
   - Store photo (Phase 2)
   - Rate card (if set up)
4. Admin verifies information:
   - Checks if store address is legitimate
   - Reviews ID documents if provided
   - Checks if rate card is reasonable
   - May contact laundry person for clarification
5. Admin makes decision:

**Option A: Approve**
- Select which societies to approve for (can approve some, reject others)
- Add optional admin notes
- Confirm approval
- Laundry person notified immediately
- Can start accepting orders from approved societies

**Option B: Reject**
- Select rejection reason:
  - Incomplete information
  - Invalid documents
  - Society quota full (too many laundry persons already)
  - Suspicious/fraudulent application
  - Other (specify)
- Add detailed explanation
- Laundry person notified with reason
- Can reapply after addressing issues

**Option C: Request More Information**
- Send message to laundry person
- Specify what's needed
- Application remains pending
- Laundry person submits additional info
- Admin reviews again

**Approval Criteria:**
- Valid phone number
- Complete address information
- Reasonable rate card pricing (not too high/low)
- Valid ID proof (Phase 2)
- Maximum laundry persons per society: 5-7 (configurable)

**Time Expectations:**
- Standard: Within 24 hours
- Urgent cases: Within 4 hours
- Weekend applications: Next business day

---

### 5.3 Manage Resident Rosters

**Flow:**

**Upload New Roster:**
1. Admin selects society
2. Downloads current roster (CSV format) for reference
3. Prepares new roster CSV with columns:
   - Phone (required, 10 digits)
   - Flat (required, format: A-404)
   - Name (optional)
4. Uploads CSV file
5. Chooses update mode:
   - Replace existing roster (full overwrite)
   - Add to existing roster (merge new entries)
6. System validates file:
   - Checks format
   - Identifies duplicates
   - Flags errors (missing phone, invalid format)
7. Admin reviews validation results:
   - Valid entries count
   - Warnings (duplicates, format issues)
   - Errors (must fix before proceeding)
8. Admin chooses:
   - Fix errors and re-upload
   - Skip invalid rows and upload valid ones
9. Confirms upload
10. System processes:
    - Imports valid entries
    - Auto-verifies residents with matching phones
    - Notifies newly verified residents
    - Updates resident database

**Manual Resident Verification:**
1. Admin views pending verification requests
2. Each request shows:
   - Resident name and phone
   - Claimed flat and society
   - Registration date
3. Admin verifies:
   - Cross-checks with society records
   - May contact resident for proof
   - May visit flat to confirm
4. Admin approves or rejects:
   - Approve: Resident can immediately book orders
   - Reject: Provide reason, resident can re-apply

**Edit Existing Roster:**
1. Admin searches for specific resident
2. Views current details
3. Can update:
   - Flat number
   - Name
   - Add/remove from roster
4. Changes take effect immediately
5. Resident notified of updates

**Roster Maintenance:**
- Regular updates recommended (quarterly)
- Remove residents who have moved out
- Add new residents promptly
- Audit roster accuracy periodically

---

### 5.4 View & Monitor Orders

**Flow:**

1. Admin accesses order monitoring dashboard
2. Views all orders in society with filters:
   - By date range
   - By status
   - By laundry person
   - By building/floor
   - Problem orders only (delays, disputes, cancellations)
3. Each order shows:
   - Order ID
   - Resident and laundry person
   - Status and timing
   - Amount
   - Any flags (delayed, disputed, unpaid)
4. Can drill into specific order for full details
5. Can take actions:
   - Contact resident or laundry person
   - Resolve disputes
   - Issue warnings
   - Escalate serious issues

**Quality Monitoring Views:**
- Orders delayed beyond SLA
- High dispute rate laundry persons
- Frequent cancellations by any party
- Payment issues (frequent non-payment)
- Rating trends (declining ratings)

---

### 5.5 Dispute Resolution

**Flow:**

1. Dispute escalated to admin (from resident issue report or laundry person dispute)
2. Admin receives notification with priority level:
   - High: Financial claim >₹500, urgent escalation
   - Medium: Financial claim ₹100-500, standard dispute
   - Low: Informational, quality feedback
3. Admin views complete dispute details:
   - Order information
   - Resident's claim with evidence (photos, description)
   - Laundry person's response with counter-evidence
   - Full order history and timeline
   - Historical data for both parties (past disputes, ratings, reliability)
4. Admin reviews evidence:
   - Photos from both sides
   - Order notes and communication
   - Pickup/delivery timestamps
   - Previous similar cases
5. Admin may:
   - Request additional information from either party
   - Contact both parties via call/message
   - Review original rate card and pricing
   - Check for patterns (serial disputer vs reliable party)
6. Admin makes decision:

**Resolution Options:**

**Full Refund to Resident:**
- Amount: Full order value
- Deducted from: Laundry person's settlement
- Reason documented
- Both parties notified

**Partial Refund to Resident:**
- Amount: Admin specifies (e.g., ₹120 of ₹240)
- Deducted from: Laundry person's settlement
- Reason documented
- Both parties notified

**No Refund (Claim Denied):**
- Resident's claim deemed invalid
- Laundry person not at fault
- No financial impact
- Reason documented
- Both parties notified

**Other Resolution:**
- Item replacement
- Re-service at no cost
- Future discount/credit
- Custom solution

7. Admin adds decision notes explaining reasoning
8. Optionally issues warning to either party:
   - Quality warning to laundry person
   - Abuse warning to resident
   - Pattern noted for future reference
9. Submits decision
10. System executes:
    - Financial transfers if applicable
    - Notifications sent to both parties
    - Dispute closed and archived
    - Metrics updated (dispute count, resolution time)

**Decision Guidelines:**

**Favor Resident When:**
- Clear photo evidence of issue
- Laundry person has history of similar issues
- Laundry person didn't respond or provided weak defense
- First-time issue for this resident

**Favor Laundry Person When:**
- Resident has pattern of frequent claims
- No evidence provided by resident
- Strong counter-evidence from laundry person
- Resident's claim inconsistent with photos

**Split Liability When:**
- Both parties partially at fault
- Evidence is inconclusive
- First occurrence for both parties
- Good faith effort from both sides

**Time Expectations:**
- High priority: Resolved within 24 hours
- Medium priority: Resolved within 48 hours
- Low priority: Resolved within 7 days

---

### 5.6 Quality Monitoring & Actions

**Flow:**

**Automated Alerts:**

System automatically flags:
- Laundry person with 3+ issues reported in 30 days
- Resident with 5+ claims in 30 days
- Laundry person with 5+ cancellations in 30 days
- Rating drop below 4.0
- Payment collection rate below 80%

**Admin Review Process:**
1. Admin receives alert with details
2. Reviews pattern and context
3. Reaches out to flagged party:
   - For laundry person: "We noticed quality concerns. Can we help?"
   - For resident: "Noticed frequent issues. Everything okay?"
4. Investigates root cause:
   - Overwork/overbooking?
   - Misunderstanding of service?
   - Genuine quality issues?
   - Abuse of system?
5. Takes appropriate action:

**For Laundry Persons:**
- Supportive conversation (most cases)
- Suggest capacity reduction if overbooked
- Training/guidance on quality
- Warning (if malicious behavior)
- Temporary suspension (serious cases)
- Permanent ban (extreme fraud/abuse)

**For Residents:**
- Educational conversation
- Clarify expectations
- Warning if abuse detected
- Temporary restriction (serious cases)
- Account suspension (extreme abuse)

**Quality Improvement Actions:**
- Rate card review and adjustment guidance
- Service standard reminders
- Best practice sharing
- Recognition of high performers

---

### 5.7 Broadcast Communications

**Flow:**

1. Admin accesses communication center
2. Selects audience:
   - All residents in society
   - All laundry persons serving society
   - Specific segment (e.g., active users only)
3. Chooses message type:
   - Announcement (service update, holiday notice)
   - Reminder (payment due, service tips)
   - Alert (urgent issue, safety concern)
4. Composes message with:
   - Subject line
   - Message body (up to 500 characters)
   - Optional image/attachment
5. Previews message
6. Schedules or sends immediately
7. Message delivered via:
   - Push notification
   - In-app notification
   - SMS (for critical alerts)

**Common Use Cases:**
- Holiday schedule changes
- New laundry person announcements
- Service quality reminders
- Payment deadline reminders
- Platform updates and features

---

## 6. Payment Processing System

### 6.1 Payment Flow Architecture

**Order Creation → Delivery → Payment → Settlement**

```
1. Order Created
   - Estimated price calculated
   - No payment collected

2. Pickup & Count Approval
   - Final price determined
   - Still no payment

3. Delivery Complete
   - Status: DELIVERED
   - Payment triggered

4. Payment Collected
   - Resident pays via UPI or Cash
   - Status: COMPLETED

5. Settlement
   - Funds allocated to laundry person
   - Platform fee deducted (if applicable)
```

---

### 6.2 UPI Payment Processing

**Flow:**

1. Resident taps "Pay via UPI" after delivery
2. System generates UPI payment request:
   - Amount: Final order amount
   - Payee UPI ID: Laundry person's registered UPI
   - Transaction note: "Order #1234 - Laundry Service"
   - Transaction ID: Unique reference
3. Deep link created for UPI apps
4. User's device opens preferred payment app (Google Pay, PhonePe, Paytm, etc.)
5. User authenticates and confirms payment in their app
6. Payment app processes transaction
7. Callback sent to our system with result
8. System waits for confirmation (30-second timeout)

**Success Scenario:**
1. Payment successful notification received
2. Transaction ID recorded
3. Order marked as "Paid via UPI"
4. Both parties notified
5. Digital receipt generated
6. Amount reflected in laundry person's settlement
7. Rating prompt shown to resident

**Failure Scenario:**
1. Payment failure notification received (or timeout)
2. User shown error with reason:
   - Insufficient balance
   - Payment declined
   - Transaction timeout
   - Network error
3. User can retry (max 3 attempts)
4. After 3 failures:
   - Suggest cash payment option
   - Option to pay later
   - Contact support option

**Technical Implementation:**
- UPI deep linking via standard UPI URI scheme
- Webhook for payment status callbacks
- Transaction reconciliation system
- Retry mechanism for failed callbacks
- Duplicate transaction prevention

**Security:**
- No storage of payment credentials
- Transaction IDs for audit trail
- Encryption for all payment data
- PCI compliance (if handling cards in future)

---

### 6.3 Cash Payment Processing

**Flow:**

1. Resident selects "Mark as Cash Paid"
2. Confirms amount paid
3. System records resident's claim
4. Status: "Payment Pending Confirmation"
5. Laundry person receives verification request
6. Laundry person confirms or disputes:

**If Confirmed:**
- Order marked as "Paid (Cash)"
- Settlement updated for laundry person
- Both parties notified
- Order completed

**If Disputed:**
- Laundry person selects "Did not receive cash"
- Automatic escalation to admin
- Order status: "Payment Disputed"
- Admin reviews case:
  - Contacts both parties
  - Reviews history (payment patterns)
  - Makes decision
- Decision executed:
  - If resident lying: Warning issued, must pay
  - If laundry person lying: Warning issued, payment recorded
  - If unclear: Admin mediates settlement

**Cash Collection Tracking:**
- Laundry person can mark "Cash Collected" during delivery
- This creates immediate confirmation
- Resident still receives notification
- Can dispute if incorrect

---

### 6.4 Payment Enforcement Policy

**First Order:**
- No restrictions
- Full trust approach
- Can pay anytime

**After 48 Hours Unpaid:**
- Soft reminder notification sent
- Restrictions applied:
  - CAN book with same laundry person (trust-based relationship)
  - CANNOT book with NEW laundry persons
  - Payment banner persistent on home screen
- No account suspension yet

**After 2+ Unpaid Orders:**
- Hard restriction applied
- Account suspended for NEW bookings
- Can still:
  - View existing orders
  - Make payments
  - Contact support
- Must clear all dues to resume booking
- Notification: "Account restricted due to unpaid orders. Pay ₹XXX to resume service."

**After 7 Days Unpaid:**
- Additional reminders (max 1 per 24 hours)
- Option for laundry person to escalate to admin
- Admin can contact resident
- May negotiate payment plan for financial hardship

**Permanent Non-Payment:**
- After 30 days with multiple attempts: Account flagged
- Collection agency option (extreme cases)
- Legal recourse available to laundry person
- Platform supports laundry person in recovery

**Rationale:**
- Balances trust with accountability
- Protects laundry persons from serial non-payers
- Allows relationship-based trust to continue
- Graduated enforcement increases compliance

---

### 6.5 Payment Disputes & Mismatches

**Scenario 1: Double Payment**
Resident claims paid cash, then also pays via UPI

**Resolution:**
1. System detects duplicate payment
2. Automatic refund of UPI payment initiated
3. Cash payment stands
4. Both parties notified
5. If intentional fraud: Warning issued

**Scenario 2: Amount Mismatch**
Resident pays ₹200 via UPI, but order amount is ₹240

**Resolution:**
1. System detects underpayment
2. Notification to resident: "Partial payment received (₹200 of ₹240)"
3. Request remaining ₹40
4. Order not completed until full payment
5. If not paid: Follows standard enforcement

**Scenario 3: Payment Not Reflected**
Resident paid via UPI but system didn't record

**Resolution:**
1. Resident reports issue
2. Provides transaction ID from payment app
3. Admin manually verifies with transaction ID
4. If verified: Manually mark as paid
5. If not found: Refund may have occurred, check with payment gateway
6. Issue resolved with admin intervention

**Scenario 4: Fraudulent Refund Request**
Resident requests refund claiming non-delivery after payment

**Resolution:**
1. Check delivery confirmation and timestamps
2. Check if photos available (Phase 2)
3. Contact laundry person for evidence
4. Review resident's history for patterns
5. Admin decides based on evidence
6. If fraud confirmed: Warning, no refund
7. If legitimate: Process refund

---

### 6.6 Offline Payment Handling

**Scenario: Resident offline during payment**

**Flow:**
1. Resident taps "Pay via UPI" but device is offline
2. System queues payment intent locally
3. Shows: "Payment will process when online"
4. When device comes online:
   - Automatic retry of payment intent
   - Up to 3 retry attempts over 24 hours
5. Notification sent when payment completes
6. If still fails: User must manually retry

**Scenario: Laundry person offline when confirming cash**

**Flow:**
1. Resident marks "Cash Paid" but LP offline
2. Verification request queued
3. LP comes online → Request delivered
4. LP confirms as normal
5. Order completes

---

## 7. Dispute Resolution System

### 7.1 Dispute Categories & Severity

**Category 1: Missing Items (High Severity)**
- Financial impact: High
- Urgency: Immediate
- Typical resolution: Refund or item return

**Category 2: Damaged Items (High Severity)**
- Financial impact: High
- Urgency: Immediate
- Typical resolution: Compensation or replacement

**Category 3: Quality Issues (Medium Severity)**
- Financial impact: Medium
- Urgency: Standard
- Typical resolution: Partial refund or re-service

**Category 4: Service Issues (Low Severity)**
- Financial impact: Low
- Urgency: Low
- Typical resolution: Apology, future credit

---

### 7.2 Dispute Scenario 1: Missing Item

**Resident Claim:**
"Blue formal shirt missing from delivery. It was included at pickup."

**Evidence Provided:**
- 2 photos (unclear if shirt visible)
- Description: "1 blue shirt missing"
- Expected: Full refund ₹240

**Laundry Person Response Option A: Accept**

**Flow:**
1. LP receives notification within 1 hour of report
2. Reviews claim
3. Decides: "Accept & Compensate"
4. Offers: ₹50 refund (for one shirt)
5. Explanation: "I sincerely apologize. I'll search for the shirt and deliver it tomorrow."
6. Resident receives offer
7. Resident choices:
   - Accept ₹50 + item return promise
   - Reject and counter with ₹100
   - Accept but escalate if item not returned
8. If accepted:
   - ₹50 deducted from LP settlement
   - Item return scheduled
   - If not returned by promised date: Auto-escalate for additional refund
9. Issue closed (pending item return)

**Laundry Person Response Option B: Dispute**

**Flow:**
1. LP responds: "Dispute this claim"
2. Counter-evidence: "All items were delivered. I have delivery photo showing all items."
3. Uploads photo from delivery
4. Explanation: "Photo clearly shows blue shirt in delivery bundle."
5. Automatic admin escalation
6. Admin reviews:
   - Resident's photos (unclear)
   - LP's delivery photo (appears to show shirt)
   - Order history: Resident's first issue, LP has 2 similar issues in past
7. Admin weighs evidence:
   - Photo from LP does show blue shirt
   - But LP has history of missing item issues
   - Resident is reliable customer
8. Admin decision: "Partial refund - ₹100"
9. Reasoning: "While delivery photo exists, pattern of missing items concerns us. Partial refund as goodwill. Warning issued to laundry person on quality."
10. Financial execution:
    - ₹100 refunded to resident (UPI)
    - ₹100 deducted from LP settlement
    - Warning recorded in LP's profile
11. Both parties notified with reasoning
12. Dispute closed

---

### 7.3 Dispute Scenario 2: Damaged Item

**Resident Claim:**
"Saree has burn mark from iron. Completely damaged."

**Evidence Provided:**
- 3 clear photos showing burn damage
- Description: "Large brown burn mark on saree pallu"
- Expected: Full refund ₹240

**Laundry Person Response: Accept & Compensate**

**Flow:**
1. LP reviews claim and photos
2. Clear damage visible in photos
3. Decides: "Accept & Compensate"
4. Offers: Full refund ₹240
5. Explanation: "I'm extremely sorry for the damage. This was a mistake during ironing. I accept full responsibility and will refund the complete amount."
6. Resident receives offer
7. Resident accepts immediately
8. Financial execution:
   - ₹240 credited to resident
   - ₹240 deducted from LP's next settlement
9. Issue closed
10. Quality incident recorded:
    - Counted toward LP's quality metrics
    - If 3rd damage incident in 30 days: Quality review triggered
    - Admin may reach out to LP for quality improvement conversation

---

### 7.4 Dispute Scenario 3: Poor Quality (Wrinkled Clothes)

**Resident Claim:**
"Clothes still wrinkled after ironing. Not acceptable quality."

**Evidence Provided:**
- 2 photos showing wrinkles
- Description: "Shirts and pants still have wrinkles"
- Expected: Re-iron at no cost

**Laundry Person Response: Partial Accept**

**Flow:**
1. LP reviews photos
2. Sees minor wrinkles (subjective quality issue)
3. Responds: "Item will be re-ironed"
4. Offers to pick up and re-iron within 24 hours at no additional cost
5. Explanation: "I apologize for not meeting your expectations. I'll pick up the items today and re-iron them properly."
6. Resident accepts offer
7. Re-service scheduled:
   - Pickup time set
   - No additional charge
   - Priority handling
8. Items re-ironed and returned
9. Resident satisfied
10. Issue closed as resolved
11. No financial impact, but noted in LP's quality record

**Alternative: Resident Rejects Re-Service**

**Flow:**
1. Resident: "I don't have time for re-service. Want partial refund."
2. LP: "Happy to offer ₹50 refund for inconvenience."
3. Resident accepts
4. ₹50 refunded
5. Issue closed

---

### 7.5 Dispute Scenario 4: Count Mismatch After Delivery

**Resident Claim:**
"I gave 12 items but only 10 returned."

**Evidence Provided:**
- Description: "2 items missing"
- No photos from pickup
- Expected: Refund for missing items

**Laundry Person Response: Strong Dispute**

**Flow:**
1. LP responds: "Dispute - All items returned"
2. Counter-evidence: "Pickup photo shows 10 items collected. Delivery photo shows same 10 items returned."
3. Uploads both photos
4. Explanation: "Resident approved count of 10 at pickup. All 10 items returned."
5. Admin escalation
6. Admin reviews:
   - Pickup count: 10 items (resident approved this)
   - Pickup photo: Shows 10 items
   - Delivery photo: Shows 10 items
   - No evidence of 12 items ever collected
7. Admin decision: "No refund - Claim invalid"
8. Reasoning: "Count was verified and approved at pickup (10 items). All photos confirm 10 items throughout. No evidence of 12 items."
9. No financial impact
10. Resident notified with explanation
11. Educational note sent: "Please verify count carefully at pickup to avoid confusion later."
12. Dispute closed in LP's favor

---

### 7.6 Dispute Scenario 5: Serial Disputer Pattern

**Background:**
Resident has filed 6 claims in last 45 days across 3 different laundry persons.

**Current Claim:**
"Missing 2 shirts - want full refund ₹300"

**Laundry Person Response:**
Disputes with delivery photo evidence.

**Admin Review:**

**Flow:**
1. Admin receives dispute
2. System flags: "High-frequency claimer"
3. Admin reviews full history:
   - 6 claims in 45 days
   - All for missing items
   - 3 different laundry persons affected
   - Varying amounts: ₹100-₹500
   - 4 out of 6 resulted in refunds
4. Admin reviews current claim:
   - LP has clear delivery photo
   - LP's first dispute in 120 days
   - LP has 4.8 rating
5. Admin decision: "Claim denied - Abuse pattern detected"
6. Reasoning: "Pattern of frequent claims suggests abuse. Delivery photo clearly shows all items. Claim denied."
7. No refund issued
8. Resident receives formal warning:
   - "Frequent claims detected. Future claims will require stronger evidence."
   - "Continued abuse may result in account restrictions."
9. Resident's future claims flagged for higher scrutiny
10. Laundry person protected from unfair claim

---

### 7.7 Dispute Resolution Decision Framework

**Evidence Hierarchy (Most to Least Weight):**

1. **Photo Evidence**
   - Timestamped photos from pickup/delivery
   - Clear, unambiguous images
   - Highest credibility

2. **System Records**
   - Count approvals with timestamps
   - GPS data (Phase 2)
   - Payment records
   - High credibility

3. **Historical Patterns**
   - Past behavior of both parties
   - Quality metrics
   - Reliability scores
   - Medium credibility

4. **Testimonials**
   - Party statements without evidence
   - Lowest credibility alone
   - Supporting evidence needed

**Decision Principles:**

**Favor Resident When:**
- Clear photo evidence of issue
- First-time issue for resident
- LP has pattern of similar issues
- LP unresponsive or provides weak defense
- Good faith evident from resident

**Favor Laundry Person When:**
- Strong counter-evidence provided
- Resident has pattern of frequent claims
- System records contradict resident's claim
- Resident's evidence weak or missing
- LP has clean track record

**Split Decision When:**
- Evidence inconclusive on both sides
- Both parties show good faith
- Minor issue with subjective standards
- First occurrence for both
- Relationship preservation important

**Financial Decisions:**

**Full Refund:**
- Clear fault with strong evidence
- High-value item damaged/missing
- LP admits full responsibility
- Safety or health concern

**Partial Refund:**
- Shared responsibility
- Inconclusive evidence
- Minor quality issue
- Goodwill gesture appropriate

**No Refund:**
- Claim clearly invalid
- Evidence contradicts claim
- Abuse pattern detected
- LP clearly not at fault

---

### 7.8 Escalation Triggers

**Automatic Admin Escalation:**

1. **Laundry Person Disputes Claim**
   - Any dispute automatically escalates
   - Cannot be resolved between parties alone

2. **Financial Claim > ₹500**
   - High-value disputes need admin review
   - Risk of significant financial impact

3. **No Response from LP (24 hours)**
   - Unresponsive LP triggers escalation
   - Protects resident from being ignored

4. **Second Issue on Same Order**
   - If additional issue reported on already disputed order
   - Indicates deeper problem

5. **High-Frequency Claimer Detected**
   - 3+ claims in 30 days
   - Automatic scrutiny applied

6. **Safety Concern Reported**
   - Theft, harassment, threatening behavior
   - Immediate admin and possibly police involvement

---

## 8. Trust & Safety Framework

### 8.1 Identity Verification

**Resident Verification:**

**Tier 1: Phone + Roster (MVP)**
- Phone OTP verification
- Society roster matching
- Admin manual verification if not in roster
- Sufficient for basic trust

**Tier 2: Enhanced (Phase 2)**
- Optional: Upload ID proof
- Optional: Add alternate contact
- Benefits: Higher trust score, faster dispute resolution

**Laundry Person Verification:**

**Tier 1: Basic (MVP)**
- Phone OTP verification
- Business address verification
- Admin approval
- Minimum to start accepting orders

**Tier 2: Enhanced (Phase 2)**
- Mandatory: Photo ID upload (Aadhaar/License/PAN)
- Mandatory: Store photo (selfie at location)
- Optional: Police verification certificate
- Benefits: "Verified" badge, higher listing priority

**Tier 3: Premium (Phase 3)**
- Police verification completed
- Business registration documents
- References from residents
- Benefits: "Premium" badge, priority support, featured listing

---

### 8.2 Rating & Review System

**Rating Collection:**

**From Residents:**
- After delivery + payment
- 1-5 stars mandatory
- Comment optional (mandatory for 1-2 stars)
- One rating per order

**Calculation:**
- Average of all ratings
- Recent ratings weighted higher (last 30 days = 70%, older = 30%)
- Minimum 5 ratings before average displayed
- Updates in real-time

**Rating Abuse Prevention:**

**For Residents:**
- 5+ consecutive 1-star ratings → Flagged for review
- Admin checks if legitimate concerns
- If abuse: Ratings discounted or removed

**For Laundry Persons:**
- Cannot rate residents (power imbalance)
- Can report problematic residents to admin
- Admin investigates and takes action if needed

**Rating Display:**
- Average rating (e.g., 4.7)
- Total orders (e.g., 127 orders)
- Distribution: 5★ (60%), 4★ (30%), 3★ (8%), 2★ (1%), 1★ (1%)
- Recent reviews (last 10) with timestamps

**Rating Impact:**

**High Rating (4.5+):**
- Higher priority in laundry person list
- "Top Rated" badge
- More resident trust and bookings

**Medium Rating (3.5-4.4):**
- Standard listing
- No special badges
- Normal visibility

**Low Rating (<3.5):**
- Lower priority in listings
- Quality review triggered
- Admin reaches out for improvement
- Risk of suspension if no improvement

---

### 8.3 Behavioral Monitoring

**Automated Flags:**

**For Laundry Persons:**
- 3+ issues reported in 30 days → Quality review
- 5+ cancellations in 30 days → Capacity check
- 10+ late deliveries in 30 days → Reliability check
- Rating drops below 4.0 → Quality intervention
- Payment collection rate < 80% → Support offered

**For Residents:**
- 5+ claims in 30 days → Abuse review
- 10+ cancellations in 30 days → Awareness reminder
- 3+ payment defaults → Credit restriction
- Multiple disputes lost → Warning issued

**Admin Actions Based on Flags:**

**Level 1: Supportive Conversation**
- Reach out to understand issues
- Offer help and guidance
- Educational resources
- No penalties

**Level 2: Warning**
- Formal notice of concerning pattern
- Explanation of potential consequences
- Offer to help improve
- Recorded in profile

**Level 3: Temporary Restriction**
- Limited functionality (e.g., max 5 orders/week)
- Enhanced monitoring
- Path to reinstatement
- Time-bound (typically 30 days)

**Level 4: Suspension**
- Account temporarily disabled
- Review and appeal process
- Requires demonstrated improvement
- Rare, for serious violations

**Level 5: Permanent Ban**
- Account permanently disabled
- Only for extreme cases (fraud, safety threats)
- Appeals reviewed by senior admin
- Very rare

---

# Part 3: Payment & Settlement Model

## 9. Payment Flow Architecture

### 9.1 Payment Philosophy

**Core Principles:**
- **No upfront payment**: Build trust, reduce friction at booking
- **Pay after delivery**: Resident verifies quality before paying
- **Flexible payment methods**: UPI and Cash to accommodate all users
- **Direct settlement**: Money flows directly to laundry person (no escrow in MVP)
- **Transparent pricing**: Rate card-based, no hidden fees
- **Fair enforcement**: Graduated restrictions for non-payment

### 9.2 Key Stakeholders

**Residents:**
- Pay for services rendered
- Expect payment flexibility
- Want transaction records

**Laundry Persons:**
- Receive payments for services
- Need reliable payment collection
- Want quick settlement visibility

**Platform:**
- Facilitates payment flow
- Ensures payment reliability
- Collects service fee (future)

**Payment Gateway:**
- Processes UPI transactions
- Handles payment confirmations
- Manages refunds

---

### 9.3 Complete Payment Journey

**Stage 1: Order Creation (No Payment)**
```
Resident books order
↓
Estimated price calculated based on rate card
↓
No payment collected
↓
Order created with status: BOOKING CREATED
```

**Stage 2: Pickup & Count Approval (Price Finalization)**
```
Laundry person picks up items
↓
Actual count entered (may differ from booking)
↓
New price calculated
↓
Resident approves count and price
↓
Final price locked
↓
Still no payment collected
```

**Stage 3: Service Delivery (Payment Trigger)**
```
Ironing completed
↓
Items delivered to resident
↓
Status: DELIVERED
↓
Payment flow triggered
↓
Resident sees payment options
```

**Stage 4: Payment Collection**
```
Resident chooses payment method:

Option A: UPI Payment
↓
Payment processed through gateway
↓
Confirmation received
↓
Order marked: PAID (UPI)

Option B: Cash Payment
↓
Resident marks "Cash Paid"
↓
Laundry person confirms receipt
↓
Order marked: PAID (CASH)
```

**Stage 5: Settlement**
```
Payment confirmed
↓
Amount allocated to laundry person
↓
Visible in their settlement dashboard
↓
Platform fee deducted (if applicable - future)
↓
Laundry person can withdraw or track
```

---

### 9.4 Payment Timing Options

**Option 1: Immediate Payment (Recommended)**
- Resident pays right after delivery
- Order completes immediately
- Best for both parties

**Option 2: Deferred Payment**
- Resident selects "Pay Later"
- 48-hour grace period before reminders
- Payment enforcements apply after grace period

**Option 3: Cash on Delivery**
- Cash exchanged during delivery itself
- Laundry person confirms receipt in app
- Immediate settlement visibility

---

## 10. Payment Collection Methods

### 10.1 UPI Payment (Primary Method)

**Technical Flow:**

1. **Payment Initiation:**
   - Resident taps "Pay via UPI"
   - System generates UPI payment request
   - Parameters:
     - Amount: Order's final price (e.g., ₹240)
     - Payee UPI ID: Laundry person's registered UPI
     - Payer UPI ID: Resident's (if saved)
     - Transaction note: "Order #1234 - Laundry Service - Clean Press"
     - Transaction reference: Unique ID generated by platform
     - Expiry: 5 minutes

2. **UPI Deep Link Creation:**
   ```
   upi://pay?
   pa=laundryperson@upi
   &pn=Clean Press Laundry
   &tn=Order 1234
   &am=240
   &cu=INR
   &tr=TXN202511081234567
   ```

3. **Payment App Launch:**
   - Deep link opens resident's preferred UPI app
   - Auto-populated payment details
   - Resident authenticates (PIN/biometric)
   - Confirms payment

4. **Payment Processing:**
   - UPI app processes transaction
   - Money transferred between bank accounts
   - Payment gateway receives status

5. **Callback & Confirmation:**
   - Payment gateway sends webhook to platform
   - Platform receives transaction status:
     - SUCCESS: Payment completed
     - FAILURE: Payment failed
     - PENDING: Awaiting confirmation
   - Timeout: 30 seconds

6. **Post-Payment Actions:**
   - Update order status to COMPLETED
   - Generate digital receipt
   - Notify both parties
   - Update laundry person's settlement
   - Trigger rating prompt

**UPI Payment Success Rate Optimization:**

- Auto-retry on timeout (1 retry)
- Clear error messages for failures
- Multiple UPI app support (GPay, PhonePe, Paytm, etc.)
- Option to change UPI app if first choice fails

**Payment Gateway Integration:**
- Razorpay, Paytm, or Cashfree
- Webhook endpoint for callbacks
- Transaction ID mapping
- Refund API access

---

### 10.2 Cash Payment

**Flow:**

**Scenario A: Cash During Delivery (Preferred)**

1. **At Delivery:**
   - Laundry person hands over items
   - Resident pays cash
   - Laundry person immediately marks in app:
     - "Cash Collected: ₹240"
   - System records timestamp

2. **Resident Notification:**
   - "Payment confirmed: ₹240 cash received by Clean Press Laundry"
   - Order status: COMPLETED
   - Rating prompt triggered

3. **Settlement Update:**
   - Amount immediately visible in laundry person's cash collection
   - No verification needed (LP confirmed at point of collection)

**Scenario B: Cash Later**

1. **After Delivery:**
   - Items delivered but cash not collected
   - Resident sees "Pay ₹240" option
   - Selects "Mark as Cash Paid"

2. **Resident Confirmation:**
   - "Did you pay ₹240 cash to Clean Press Laundry?"
   - Resident confirms

3. **Verification Request:**
   - Laundry person receives: "Resident claims cash payment of ₹240. Confirm?"
   - Order status: "Payment Pending Confirmation"

4. **Laundry Person Response:**
   - **If Confirms:** Order marked COMPLETED, settlement updated
   - **If Disputes:** "Did not receive cash" → Escalates to admin

**Cash Payment Dispute Resolution:**

Admin reviews:
- Resident's payment claim timestamp
- Laundry person's dispute timestamp
- Historical payment patterns for both parties
- Communication logs between parties

Typical outcomes:
- If resident has reliable payment history: Benefit of doubt to resident
- If LP has clean record: Investigate resident's claim more closely
- If unclear: Admin contacts both parties for clarification

---

### 10.3 Payment Split Scenario (Phase 2)

**Use Case:** Resident wants to pay partially via UPI and partially via cash

**Example:**
- Order total: ₹500
- UPI payment: ₹300
- Cash payment: ₹200

**Flow:**

1. Resident selects "Split Payment"
2. Specifies UPI amount: ₹300
3. Completes UPI transaction
4. Marks remaining ₹200 as cash
5. Laundry person confirms ₹200 cash receipt
6. Order marked as fully paid
7. Settlement shows: ₹300 (UPI) + ₹200 (Cash)

---

## 11. Settlement Model

### 11.1 Laundry Person Settlement Dashboard

**Real-Time Visibility:**

Laundry person sees:

**This Week Summary:**
```
Total Earnings: ₹12,450
├─ Collected: ₹9,800 (78.7%)
│  ├─ UPI Received: ₹6,200
│  └─ Cash Collected: ₹3,600
└─ Pending: ₹2,650 (21.3%)
   ├─ Pending UPI: ₹1,850
   └─ Pending Cash: ₹800
```

**Pending Payments Breakdown:**
```
Order #1234 | Raj Kumar (A-404)
Amount: ₹240
Status: Unpaid
Delivered: 2 days ago
Actions: [Mark Cash Received] [Send Reminder]

Order #1235 | Priya Shah (B-201)
Amount: ₹180
Status: Awaiting Payment
Delivered: Today
Actions: [Send Reminder]

Order #1236 | Amit Verma (C-305)
Amount: ₹320
Status: Payment Pending Confirmation
Delivered: Yesterday
Note: Resident claimed cash paid, awaiting your confirmation
Actions: [Confirm Receipt] [Dispute]
```

**Historical View:**
- Last 7 days, 30 days, 90 days, All time
- Filter by: Society, Payment method, Status
- Export CSV for accounting

---

### 11.2 Settlement Flow for UPI Payments

**Direct Settlement (MVP Model):**

```
Resident pays ₹240 via UPI
↓
Payment goes directly to laundry person's UPI ID
↓
Laundry person's bank account credited (instant)
↓
Platform records transaction
↓
Laundry person sees "₹240 received via UPI" in settlement
↓
No withdrawal needed (already in their account)
```

**Advantages:**
- Instant settlement for laundry person
- No platform escrow needed
- Simple implementation
- Lower platform liability

**Disadvantages:**
- Platform cannot hold funds for disputes
- Refunds require requesting money back from laundry person
- No automatic platform fee deduction

---

### 11.3 Settlement Flow for Cash Payments

**Self-Reported Model:**

```
Laundry person collects ₹240 cash
↓
Marks in app: "Cash Collected ₹240"
↓
Platform records in settlement
↓
Physical cash already with laundry person
↓
Settlement shows: "₹240 cash collected"
↓
No transfer or withdrawal needed
```

**Verification:**
- Resident receives notification of cash collection claim
- Can dispute if incorrect
- Disputes escalate to admin

---

### 11.4 Alternative Settlement Model: Platform Escrow (Phase 2)

**For UPI Payments Only:**

```
Resident pays ₹240 via UPI
↓
Payment goes to platform's escrow account
↓
Funds held for 24-48 hours (dispute window)
↓
If no dispute: Funds released to laundry person
↓
Laundry person withdrawal options:
  - Auto-transfer to bank (daily/weekly)
  - Manual withdrawal request
  - Minimum threshold (e.g., ₹500)
```

**Advantages:**
- Platform controls refunds
- Easier dispute resolution
- Automatic platform fee deduction
- Fraud protection

**Disadvantages:**
- Delayed settlement for laundry person
- Requires payment gateway accounts
- More complex implementation
- Regulatory compliance (KYC, etc.)

---

## 12. Platform Revenue Model

### 12.1 Society Subscription Model ⭐ ADOPTED

**Core Philosophy:**
- Societies pay a monthly subscription fee
- Zero commission on orders
- Laundry persons keep 100% of their earnings
- Unlimited orders within the society
- Platform revenue is predictable and scalable

---

### 12.2 Pricing Tiers

**Based on Society Size:**

```
STARTER (100-300 flats)
Monthly Fee: ₹5,000
Features:
- Unlimited orders
- Up to 3 approved laundry persons
- Basic analytics
- Email support
- Order & payment tracking
- Dispute resolution

GROWTH (301-600 flats)
Monthly Fee: ₹10,000
Features:
- Unlimited orders
- Up to 5 approved laundry persons
- Advanced analytics
- Priority support
- Custom rate card templates
- Monthly performance reports

ENTERPRISE (601+ flats)
Monthly Fee: ₹20,000
Features:
- Unlimited orders
- Unlimited laundry persons
- Premium analytics & insights
- Dedicated account manager
- White-label customization (Phase 3)
- API access (Phase 3)
- Custom integrations
```

---

### 12.3 Why Society Subscription Works Better

**Comparison with Commission Model:**

| Factor | Commission Model (5%) | Society Subscription |
|--------|----------------------|---------------------|
| **Platform Revenue** | Variable, unpredictable | Fixed, predictable ✅ |
| **Laundry Person Earnings** | 95% of order value | 100% of order value ✅ |
| **Resident Costs** | Higher (fees passed on) | Lower (no markup) ✅ |
| **Sales Complexity** | High (explain fees) | Low (simple pricing) ✅ |
| **Break-even Speed** | 18-24 months | 6-9 months ✅ |
| **Stakeholder Alignment** | Conflicting interests | Aligned incentives ✅ |
| **Pricing Transparency** | Hidden in order cost | Upfront and clear ✅ |
| **Society Buy-in** | Resistance (resident cost) | Strong (quality service) ✅ |

---

### 12.4 Value Proposition to Societies

**What Society Management Gets:**

**1. Resident Satisfaction (₹10,000/month investment)**
- Modern, convenient service
- Verified service providers
- Quality monitoring & ratings
- Digital payment tracking
- Reduced complaints to RWA office

**2. Administrative Efficiency**
- Automated vendor management
- Digital dispute resolution
- No cash handling issues
- Transparent rate cards
- Performance analytics

**3. Quality Assurance**
- Background-verified laundry persons
- Rating system enforcement
- Issue tracking & resolution
- Service level monitoring
- Platform handles all problems

**4. ROI Calculation**
```
Cost: ₹10,000/month
Benefits:
- Reduced admin time: ~10 hours/month (₹5,000 value)
- Fewer complaints: ~20 calls/month avoided
- Resident satisfaction: Priceless
- Professional image: Modern society amenity

Easily justifiable from society maintenance fund
```

---

### 12.5 Revenue Projections

**Conservative Scenario (12 months):**

```
Month 3:   5 societies  × ₹7,000 avg  = ₹35,000/month
Month 6:   15 societies × ₹8,000 avg  = ₹1,20,000/month
Month 9:   30 societies × ₹9,000 avg  = ₹2,70,000/month
Month 12:  50 societies × ₹10,000 avg = ₹5,00,000/month

Year 1 Total Revenue: ~₹25,00,000
Operating Costs: ~₹18,00,000
Year 1 Net: ₹7,00,000 profit
```

**Aggressive Scenario (12 months):**

```
Month 3:   10 societies × ₹8,000 avg  = ₹80,000/month
Month 6:   30 societies × ₹10,000 avg = ₹3,00,000/month
Month 9:   60 societies × ₹12,000 avg = ₹7,20,000/month
Month 12:  100 societies × ₹12,000 avg = ₹12,00,000/month

Year 1 Total Revenue: ~₹60,00,000
Operating Costs: ~₹25,00,000
Year 1 Net: ₹35,00,000 profit
```

---

### 12.6 Go-to-Market Strategy

**Phase 1: Pilot (Month 1-3)**
```
Target: 5 societies
Offer: 3 months free
Goal: Prove value, get testimonials
Investment: ₹0 revenue, focus on learning
```

**Phase 2: Early Adopter (Month 4-9)**
```
Target: 20 societies
Offer: 50% discount (₹5k for Growth tier)
Goal: Scale and refine
Revenue: ₹1,00,000 - ₹3,00,000/month
```

**Phase 3: Standard Pricing (Month 10+)**
```
Target: All new societies
Offer: Full pricing with 14-day trial
Goal: Sustainable growth
Revenue: ₹5,00,000+/month by Month 12
```

---

### 12.7 Why No Commission Model

**Problems with Commission/Transaction Fees:**

1. **Misaligned Incentives**
   - Platform wants more orders (volume)
   - Quality may suffer from pressure
   - Laundry persons feel squeezed

2. **Hidden Costs**
   - Residents pay more (fees passed on)
   - Pricing becomes opaque
   - Trust issues

3. **Variable Revenue**
   - Hard to forecast
   - Difficult to raise funding
   - Can't plan hiring/scaling

4. **Competitive Disadvantage**
   - Other platforms charge 10-20%
   - Race to bottom on pricing
   - Unsustainable long-term

**Society Subscription Solves All These:**
- Fixed, predictable revenue
- No conflict of interest
- Laundry persons keep 100%
- Lower costs for residents
- Easy to explain and sell
- Aligns everyone's interests

---

### 12.8 Implementation Notes

**Billing System:**

**Monthly Billing Cycle:**
```
1. Society onboarded on 15th January
2. First bill generated: 1st February
3. Payment due: 5th February
4. Grace period: 3 days
5. Service paused if unpaid: 8th February
6. Auto-resume on payment
```

**Payment Methods for Societies:**
- Bank transfer (NEFT/RTGS)
- UPI (for smaller societies)
- Standing instruction (auto-debit)
- Cheque (if needed)

**Contract Terms:**
- Quarterly commitment (3 months minimum)
- Annual discount: 10% off (pay ₹1,08,000 instead of ₹1,20,000)
- No setup fees
- Cancel anytime after minimum period

---

### 12.9 Future Revenue Streams (Phase 2-3)

**Additional Services (Optional Add-ons):**

```
Premium Features:
├─ Advanced Analytics Dashboard: +₹2,000/month
├─ WhatsApp Integration: +₹1,000/month
├─ Custom Branding: +₹3,000/month
├─ Priority Support (24/7): +₹2,000/month
└─ Insurance Coverage: +₹5,000/month

Resident can pay for these too
```

**Marketplace Revenue (Phase 3):**
- Other home services (plumber, electrician)
- Advertise to verified residents
- Commission on other services: 10-15%

---

### 12.10 Comparison Summary

**Final Decision: Society Subscription Model**

**Reasons:**
✅ Predictable revenue (crucial for sustainability)
✅ Faster path to profitability (6-9 months vs 18-24 months)
✅ Zero conflict of interest (everyone wins)
✅ Easy to explain and sell to societies
✅ Laundry persons keep 100% (happy providers = better service)
✅ Lower costs for residents (no hidden fees)
✅ Scalable and sustainable business model
✅ Aligns with "trust and leniency" philosophy
✅ Easier to raise funding (predictable MRR)
✅ Professional, B2B credibility

**No commission or transaction fees on orders.**
**Platform is funded entirely by society subscriptions.**

---

### 12.2 Fee Implementation (When Introduced)

**Transparent Fee Display:**

At booking:
```
Items: 12 pieces
Rate card total: ₹240
Platform fee (10%): ₹24
─────────────────────
Total to pay: ₹264
```

OR (if platform absorbs fee):
```
Items: 12 pieces
Total: ₹240
(Includes small platform fee)
```

**Settlement Display:**

Laundry person sees:
```
Order #1234: ₹240
Platform fee (10%): -₹24
Net earnings: ₹216
```

---

## 13. Payment Enforcement & Collection

### 13.1 Graduated Enforcement Policy

**Timeline:**

**Day 0: Delivery Completed**
- Payment option presented
- No restrictions yet
- Soft CTA: "Pay Now" or "Pay Later"

**Day 1-2: Grace Period**
- No reminders
- No restrictions
- Building trust

**Day 3: First Reminder**
- Gentle notification: "Reminder: Order #1234 payment pending (₹240)"
- Still no restrictions
- One-tap payment link

**Day 5: Restriction Applied**
- Cannot book NEW laundry persons
- Can still book SAME laundry person (trust relationship)
- Banner: "Pay ₹240 to unlock all services"

**Day 7: Escalation**
- Additional reminder
- Laundry person can escalate to admin
- Admin may reach out to resident

**Day 14: Hard Restriction (If 2+ Unpaid Orders)**
- Account suspended for all bookings
- Must clear dues to resume
- Can view orders and make payments
- Contact support for payment issues

**Day 30: Collection Action**
- Admin involvement
- May negotiate payment plan
- Possible external collection (extreme cases)

---

### 13.2 Payment Reminder Strategy

**Frequency Limits:**
- Maximum 1 reminder per 24 hours per order
- Total maximum 5 reminders per order
- After 5 reminders: Admin escalation

**Reminder Channels:**
- Day 3: Push notification
- Day 5: Push + SMS
- Day 7: Push + SMS + In-app banner
- Day 10: Email (if provided)
- Day 14: Admin call/message

**Reminder Tone:**
- Days 3-5: Friendly reminder
- Days 7-10: Urgent reminder
- Days 14+: Formal notice

---

### 13.3 Incentives for Timely Payment

**Positive Reinforcement (Phase 2):**

**On-Time Payment Streak:**
- Pay within 24 hours for 5 consecutive orders → "Reliable Customer" badge
- Benefits:
  - Priority support
  - Early access to new features
  - Possible discounts from laundry persons

**Instant Payment Benefits:**
- Pay immediately after delivery
- Get ₹10 platform credit (usable on next order)
- Limited-time promotion to encourage adoption

**Loyalty Points:**
- Earn 1 point per ₹100 paid on time
- Redeem 100 points = ₹50 discount
- Encourages repeat usage and timely payment

---

### 13.4 Non-Payment Recovery

**For Laundry Persons:**

**Day 7 Options:**
- Send payment reminder (platform-facilitated)
- Escalate to admin for follow-up
- Platform contacts resident on LP's behalf

**Day 14 Options:**
- Admin attempts resolution
- May negotiate payment plan
- Document non-payment

**Day 30 Options:**
- Admin makes final attempt
- If persistent non-payment:
  - Resident account flagged/suspended
  - LP can choose to pursue independently
  - Platform provides order records for legal recourse (if needed)

**Protection for Laundry Persons:**
- Platform supports collection efforts
- Provides transaction records
- Admin mediation available
- Serious cases: Legal support guidance

---

## 14. Financial Reconciliation

### 14.1 Daily Reconciliation

**For Platform:**

Daily automated checks:
```
UPI Payments Received:
- Gateway reports: ₹125,000
- Platform records: ₹125,000
- Match: ✓

Cash Payments Claimed:
- Laundry persons claimed: ₹85,000
- Resident confirmations: ₹82,000
- Discrepancy: ₹3,000 (under verification)

Refunds Processed:
- Total refunds: ₹2,400
- Deducted from LP settlements: ₹2,400
- Match: ✓
```

**Discrepancy Resolution:**
- UPI mismatch: Check gateway logs, retry callbacks
- Cash mismatch: Review disputed cash payments
- Refund mismatch: Audit dispute resolutions

---

### 14.2 Settlement Reconciliation for Laundry Persons

**Weekly Statement:**
```
Settlement Summary: Nov 1-7, 2025
Clean Press Laundry

EARNINGS
Total orders completed: 45
Total order value: ₹9,450

COLLECTIONS
UPI payments: ₹6,200 (65.6%)
  - Settled to your account: ₹6,200
Cash payments: ₹2,800 (29.6%)
  - Collected by you: ₹2,800
Pending: ₹450 (4.8%)
  - 3 orders unpaid

DEDUCTIONS
Refunds issued: ₹240 (1 order)
Platform fee: ₹0 (current promotion)
Total deductions: ₹240

NET EARNINGS
Total: ₹9,210
Already received: ₹9,000 (UPI + Cash)
Pending: ₹450
Adjustment: -₹240 (refunds)
```

---

### 14.3 Audit Trail

**Every Financial Transaction Recorded:**

**UPI Payment:**
```
Transaction ID: TXN202511081234567
Order ID: #1234
Timestamp: 2025-11-08 18:30:45
Amount: ₹240
Method: UPI
Payer: Resident (9876543210)
Payee: Laundry Person (9123456780)
Gateway: Razorpay
Gateway TXN ID: pay_AbC123XyZ
Status: SUCCESS
```

**Cash Payment:**
```
Transaction ID: CASH202511081234568
Order ID: #1235
Timestamp: 2025-11-08 14:20:15
Amount: ₹180
Method: CASH
Claimed by: Laundry Person (9123456780)
Confirmed by: Resident (9876543210)
Status: CONFIRMED
```

**Refund:**
```
Transaction ID: REF202511081234569
Order ID: #1234
Timestamp: 2025-11-09 10:15:30
Amount: ₹120 (partial refund)
Reason: Missing item dispute
Approved by: Admin
Original payment: TXN202511081234567
Deducted from: LP settlement
Status: PROCESSED
```

---

## 15. Edge Cases & Scenarios

### 15.1 Double Payment

**Scenario:**
Resident pays ₹240 via UPI, then also marks "Cash Paid" for same order.

**Detection:**
- System detects two payment attempts on same order
- Flags for review

**Resolution:**
1. Most recent payment cancelled automatically
2. If UPI already processed: Refund initiated
3. If cash also claimed: Admin investigates
4. Resident notified: "Duplicate payment detected. ₹240 refunded."
5. Order marked as paid once (whichever came first)

---

### 15.2 Partial Payment

**Scenario:**
Order total ₹240, resident pays ₹200 via UPI.

**Detection:**
- System checks payment amount vs order amount
- Mismatch detected: ₹40 short

**Resolution:**
1. Order marked as "Partially Paid"
2. Notification: "Partial payment received (₹200 of ₹240). Pay remaining ₹40."
3. Resident can pay remaining via UPI or cash
4. Order completes only after full payment

---

### 15.3 Overpayment

**Scenario:**
Order total ₹240, resident pays ₹300 via UPI.

**Detection:**
- System checks payment amount vs order amount
- Overpayment detected: ₹60 extra

**Resolution:**
1. Order marked as "Paid (Overpaid)"
2. Options presented:
   - Refund ₹60 to resident
   - Apply ₹60 as credit to next order
   - Tip laundry person ₹60
3. Resident chooses
4. Action executed accordingly

---

### 15.4 Payment Not Reflected

**Scenario:**
Resident paid via UPI but system didn't record it.

**Detection:**
- Resident reports: "I paid but order still shows unpaid"
- Provides UPI transaction ID from their payment app

**Resolution:**
1. Admin manually verifies with payment gateway
2. Looks up transaction ID
3. If found and successful:
   - Manually mark order as paid
   - Update settlement
   - Apologize for delay
4. If not found:
   - Check if payment actually debited from resident's account
   - If debited but not received: Gateway issue, escalate
   - If not debited: Ask resident to retry payment

---

### 15.5 Fraudulent Refund Request

**Scenario:**
Resident requests refund claiming non-delivery, but laundry person has proof of delivery.

**Detection:**
- Laundry person disputes refund request
- Provides delivery photo with timestamp
- GPS data shows presence at delivery location (Phase 2)

**Resolution:**
1. Admin reviews evidence
2. Delivery proof is clear
3. Refund denied
4. Resident warned about false claims
5. Pattern tracked for future abuse detection

---

### 15.6 Laundry Person Account Closure with Pending Payments

**Scenario:**
Laundry person wants to close account but has ₹2,500 in pending payments from residents.

**Resolution:**
1. Platform prevents account closure until:
   - All pending payments collected OR
   - Laundry person waives pending amounts
2. Options provided:
   - Wait for payments to clear
   - Send final payment reminders
   - Admin helps with collection
   - Accept loss and waive amounts
3. Once resolved, account can be closed

---

### 15.7 Platform Fee Refund (When Fees Introduced)

**Scenario:**
Order refunded due to dispute. Platform fee was 10% (₹24 of ₹240).

**Policy Options:**

**Option A: Refund Full Amount (Customer-Friendly)**
```
Resident paid: ₹264 (inc. ₹24 platform fee)
Order refunded: ₹264
Platform absorbs: ₹24 loss
Laundry person deduction: ₹240
```

**Option B: Refund Minus Platform Fee**
```
Resident paid: ₹264
Order refunded: ₹240 (excluding platform fee)
Platform keeps: ₹24
Laundry person deduction: ₹240
```

**Recommendation:** Option A for customer satisfaction, Option B for platform sustainability (context-dependent)

---

### 15.8 Bulk Payment Scenario (Phase 3)

**Scenario:**
Resident has 5 unpaid orders totaling ₹1,200 with same laundry person.

**Flow:**
1. Resident accesses "Pending Payments"
2. Selects multiple orders
3. Sees total: ₹1,200
4. Option: "Pay All via UPI"
5. Single UPI transaction for ₹1,200
6. Platform distributes payment across 5 orders
7. All orders marked as paid
8. Single receipt with breakdown

---

### 15.9 Payment Plan for Large Amount (Phase 3)

**Scenario:**
Resident owes ₹5,000 (multiple unpaid orders) but facing financial hardship.

**Flow:**
1. Resident contacts admin explaining situation
2. Admin reviews:
   - Total owed
   - Resident's history
   - Laundry person's position
3. Admin proposes payment plan:
   - Pay ₹2,000 immediately
   - Pay ₹1,500 in 15 days
   - Pay ₹1,500 in 30 days
4. Resident agrees
5. Laundry person agrees (goodwill)
6. Payment plan activated
7. Reminders sent at each milestone
8. Account restrictions lifted after first payment

---

## 16. Summary & Recommendations

### 16.1 MVP Payment Model

**For Launch:**
- No upfront payment
- Direct UPI settlement to laundry person
- Cash payment self-reported with verification
- No platform fees (0% commission)
- Graduated payment enforcement
- Basic dispute resolution with manual refunds

**Why This Works:**
- Simple implementation
- Fast go-to-market
- Builds trust with zero fees
- Minimal regulatory requirements
- Laundry persons get instant settlement

---

### 16.2 Phase 2 Enhancements

- Platform escrow for UPI (24-48 hour hold)
- Automated refund processing
- Split payment options
- Bulk payment for multiple orders
- Enhanced fraud detection

---

### 16.3 Phase 3 Revenue Generation

- Introduce platform commission (5-10%)
- Subscription plans for laundry persons
- Premium features (analytics, priority listing)
- Value-added services (insurance, financing)

---

### 16.4 Key Metrics to Track

**Payment Health:**
- Payment collection rate (target: 95%+)
- Average time to payment (target: <24 hours)
- Payment method split (UPI vs Cash)
- Default rate (target: <5%)

**Settlement Health:**
- Settlement accuracy (target: 99.9%+)
- Reconciliation discrepancies
- Refund processing time (target: <48 hours)

**User Satisfaction:**
- Payment ease rating
- Dispute resolution satisfaction
- Settlement clarity feedback

---

**End of Complete PRD Document**
