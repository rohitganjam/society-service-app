# App Functionality Summary

**Version:** 1.0
**Date:** November 17, 2025
**Purpose:** Complete overview of app functionality across all user personas

---

## 🔄 **Multi-Society Context Switching**

Users can register and manage multiple residences across different societies (e.g., primary home + vacation home).

**Capabilities:**
- Register in multiple societies with different units
- Switch active society context to access location-specific vendors and services
- View all residences with verification status (verified, pending, rejected)
- See which society is currently active via persistent header
- Orders are tied to active society to prevent wrong-location mistakes

**Database (Implemented):**
- `is_primary` - Main residence (ONE per user, trigger enforced)
- `is_active` - Current society context (ONE per user, trigger enforced)
- Functions: `set_active_society()`, `get_user_active_society()`, `get_user_all_residences()`

**Mobile UX:**
- **Active society header:** Always visible, shows current society and unit (e.g., "🏢 Maple Gardens, Flat A-404")
- **Society switcher:** Tap header to open bottom sheet with all residences
- **Society cards:** Show status badges (✓ Active, "Switch →" button for verified, "⏳ Pending" for unverified)
- **Context-aware orders:** Order creation shows "📍 For: [Active Society]" to prevent mistakes

**API Endpoints (Implemented):**
- `GET /api/v1/residents/{user_id}/residences` - List all residences
- `POST /api/v1/residents/{user_id}/switch-society` - Switch active context

---

## 🏠 **Resident (Customer)**

**What they do:**
- Browse service categories: Laundry (live), Vehicle Services, Home Services, Personal Care (coming soon)
- Within each category, browse service providers and their offerings
- **Vendor Filtering & Discovery:**
  - **Default view:** See vendors assigned to their building/block/phase (filtered automatically)
  - **Option to view all:** Can manually choose to see all society vendors if needed
  - **Use cases for viewing all vendors:**
    - Emergency situations requiring immediate service
    - Specific vendor preferences or familiarity
    - Higher quality or specialty services from other areas
- **Create separate orders per category** (one laundry order, one vehicle order, etc.)
- **Within each order**, can mix multiple service types from that category:
  - Laundry order: ironing + washing + dry cleaning items together
  - Vehicle order: car wash + bike wash + detailing together
  - Home order: gardening + plumbing tasks together
- Select pickup/service time and address
- Track order status in real-time (each service type has its own completion workflow)
- Approve item count changes if vendor finds discrepancies
- Pay via UPI or mark cash payment after delivery
- Rate and review service providers per service type
- Report issues/disputes if needed

**Key flow:** Browse categories → Filter vendors (default: assigned) → Select provider → Add items to cart (within category) → Schedule pickup → Track per service → Pay → Review

**Example**:
- Creates **Laundry Order #123**: 5 shirts (ironing) + 2 suits (dry cleaning) = ₹350
- Creates **Vehicle Order #124**: 1 car wash + 1 bike wash = ₹400
- Two separate orders, each with mixed service types within their category

---

## 👔 **Service Provider (Vendor)**

**What they do:**
- Register and select which service categories/types they offer:
  - **Laundry Services**: Ironing, Washing, Dry Cleaning, Washing Only
  - **Vehicle Services**: Car Wash, Bike Wash, Detailing, Interior Cleaning
  - **Home Services**: Gardening, Plumbing, Electrical, Pest Control
  - **Personal Care**: Barber, Salon, Spa services
- Can offer services across multiple categories (e.g., runs laundry + car wash business)
- **Service Area Assignment:**
  - Society admin assigns vendor to specific service areas (buildings, blocks, phases, or entire society)
  - Assignment serves as default for which residents see this vendor
  - Residents can still choose to order from vendor outside default area
- Set up separate rate cards with pricing for each service type
- **Define completion workflow per service type**:
  - **Laundry/Ironing**: Pickup → Count → Iron → Ready → Deliver
  - **Laundry/Dry Cleaning**: Pickup → Count → Dry Clean → Quality Check → Ready → Deliver
  - **Vehicle/Car Wash**: Schedule → Arrive → Wash → Vacuum → Polish → Complete
  - **Home/Gardening**: Schedule → Arrive → Trim Plants → Mow Lawn → Clean → Complete
  - **Home/Plumbing**: Schedule → Arrive → Diagnose → Fix → Test → Complete
- **Order Management:**
  - **Default view:** See all assigned requests from their designated service areas
  - **Filter option:** Can filter by building/phase/service type if needed
  - Receive order notifications grouped by category and service type
- View dashboard with today's tasks organized by category
- Update status independently for each service type within an order
- Complete each service according to its specific workflow
- Track earnings and settlements by category and service type
- Respond to customer disputes

**Key flow:** Register → Select categories → Get assigned to service areas → Setup rate cards → Define workflows → Receive orders (from assigned areas) → Execute per service workflow → Get paid

**Example workflow tracking**:
```
Laundry Order #123 (Mixed):
├─ Ironing (5 shirts)
│   ├─ ✅ Picked up (Day 1)
│   ├─ ✅ Counted (Day 1)
│   ├─ ✅ Ironed (Day 2)
│   ├─ ✅ Ready for delivery (Day 2)
│   └─ ⏳ Delivered (Pending)
│
└─ Dry Cleaning (2 suits)
    ├─ ✅ Picked up (Day 1)
    ├─ ✅ Counted (Day 1)
    ├─ ✅ Dry cleaned (Day 3)
    ├─ ✅ Quality checked (Day 4)
    ├─ ✅ Ready for delivery (Day 5)
    └─ ⏳ Delivered (Pending)

Vehicle Order #124 (Mixed):
├─ Car Wash
│   ├─ ✅ Scheduled (10 AM)
│   ├─ ✅ Arrived on site
│   ├─ ✅ Exterior washed
│   ├─ ✅ Interior vacuumed
│   ├─ ✅ Polished
│   └─ ✅ Complete
│
└─ Bike Wash
    ├─ ✅ Scheduled (10:30 AM)
    ├─ ✅ Washed
    ├─ ✅ Dried
    └─ ✅ Complete
```

---

## 🏢 **Society Admin**

**What they do:**
- Approve/reject vendor registrations for their society (across all categories)
- **Assign vendors to service areas** within the society:
  - Assign to **entire society** (all buildings/blocks and phases)
  - Assign to **specific building(s)/block(s)**
  - Assign to **specific phase(s)** (groups of households within layouts)
  - Can assign one vendor to multiple service areas
- Upload resident rosters (phone numbers + flat numbers) for instant verification
- Monitor all orders in their society (laundry, vehicle, home services - all categories)
- View order completion by service type workflow
- Resolve escalated disputes between residents and vendors
- Manage society subscription billing (same fee covers all service categories)
- View analytics: completion rates, average time per service type workflow
- Track vendor performance by category and service type

**Key flow:** Approve vendors → Assign to service areas → Manage rosters → Monitor activity → Resolve disputes

**Vendor Assignment Examples:**
- Vendor A: Assigned to entire society (serves all residents)
- Vendor B: Assigned to Building 1 and Building 2 only
- Vendor C: Assigned to Phase 1 households in independent house layout
- Vendor D: Assigned to Building 3, Floor 1-5 only

---

## 💼 **Super Admin (Platform)**

**What they do:**
- Manage multiple societies and their subscriptions
- Generate and track subscription invoices
- Handle overdue payments (suspend societies if needed)
- Monitor platform-wide metrics across all service categories
- **Manage service categories and workflows**:
  - Add new parent categories (Laundry, Vehicle, Home, etc.)
  - Add service subcategories within each parent
  - **Define workflow steps per service type**:
    ```
    Ironing: 5 steps (Pickup → Count → Iron → Ready → Deliver)
    Dry Cleaning: 6 steps (Pickup → Count → Dry Clean → QC → Ready → Deliver)
    Car Wash: 6 steps (Schedule → Arrive → Wash → Vacuum → Polish → Complete)
    Gardening: 6 steps (Schedule → Arrive → Trim → Mow → Clean → Complete)
    Plumbing: 6 steps (Schedule → Arrive → Diagnose → Fix → Test → Complete)
    ```
  - Set default turnaround times per service
  - Configure pricing models (per-item, per-service, hourly)
- Activate/deactivate categories when ready to launch
- Handle critical escalations
- View completion metrics per workflow step
- Analyze bottlenecks in service workflows

**Key flow:** Onboard societies → Configure categories/workflows → Manage billing → Activate categories → Monitor health

---

## 🏗️ **Society Organizational Structure**

### Generic Hierarchical Model

All societies use a **flexible, self-defining hierarchical structure** stored in a single table. Society admins can create any structure that matches their layout without requiring schema changes.

#### Core Concept: Nodes in a Tree

Every part of a society (from the society itself down to individual units) is represented as a **node** in a tree structure:

```
Society (Root Node)
├── Node (Building/Phase/Wing/Section/etc.)
│   ├── Node (Floor/Sub-section/etc.) [Optional depth]
│   │   ├── Node (Unit: Flat/House/Villa)
│   │   │   └── Node (Floor within unit) [Optional for multi-story units]
│   │   └── Node (Unit)
│   └── Node (Unit) [If no intermediate floors]
└── Node (Building/Phase/etc.)
```

**Key Features:**
- ✅ **Unlimited Depth:** Add as many levels as needed
- ✅ **Flexible Naming:** Use any terminology (Building, Tower, Phase, Wing, etc.)
- ✅ **No Schema Changes:** New structures = insert nodes, not alter tables
- ✅ **Path-Based Queries:** Efficient ancestor/descendant lookups

---

### Structure Examples

#### **1. Apartment Complex**
```
Green Valley Apartments (Society Node)
├── Building A (Node: BUILDING)
│   ├── Floor 1 (Node: FLOOR)
│   │   ├── Flat A-101 (Node: UNIT)
│   │   └── Flat A-102 (Node: UNIT)
│   └── Floor 2 (Node: FLOOR)
│       ├── Flat A-201 (Node: UNIT)
│       └── Flat A-202 (Node: UNIT)
└── Tower B (Node: TOWER)
    ├── Flat B-101 (Node: UNIT) [No intermediate floors]
    └── Flat B-102 (Node: UNIT)
```

**Path Examples:**
- Flat A-101 path: `1.2.4.6` (Society → Building A → Floor 1 → Flat A-101)
- Flat B-101 path: `1.3.8` (Society → Tower B → Flat B-101)

---

#### **2. Independent House Layout**
```
Sunrise Villas (Society Node)
├── Phase 1 (Node: PHASE)
│   ├── House 101 (Node: UNIT)
│   │   ├── Ground Floor (Node: FLOOR) [Multi-story house]
│   │   └── First Floor (Node: FLOOR)
│   └── House 102 (Node: UNIT) [Single-story]
└── Phase 2 (Node: PHASE)
    └── House 201 (Node: UNIT)
```

**Path Examples:**
- House 101 Ground Floor: `10.11.13.15` (Society → Phase 1 → House 101 → GF)
- House 102: `10.11.14` (Society → Phase 1 → House 102)

---

#### **3. Mixed/Custom Structure**
```
Metro Heights (Society Node)
├── North Wing (Node: WING)
│   └── Flat NW-101 (Node: UNIT)
├── Garden Villas Section (Node: SECTION)
│   ├── Villa 1 (Node: UNIT)
│   └── Villa 2 (Node: UNIT)
└── Commercial Block (Node: BLOCK)
    └── Shop C-01 (Node: UNIT)
```

**Node Types Supported:**
- `SOCIETY`, `BUILDING`, `TOWER`, `BLOCK`, `WING` (apartments)
- `PHASE`, `SECTION`, `ZONE` (layouts)
- `FLOOR`, `UNIT` (universal)
- **Custom types:** Society admins can define any node type

---

### Vendor Assignment by Hierarchy

Vendors can be assigned to **any node** in the hierarchy. Assignment automatically includes all descendant nodes.

#### **Assignment Flexibility:**

**Level 1 - Society-Wide:**
```
"Premium Services"
├── Assigned to: Society Root (entire society)
└── Serves: ALL units in every building/phase
```

**Level 2 - Building/Phase:**
```
"QuickWash Laundry"
├── Assigned to: Building A, Tower B
└── Serves: All floors and units in Building A and Tower B
```

**Level 3 - Floor-Specific:**
```
"Luxury Cleaners"
├── Assigned to: Floor 14, Floor 15 (premium floors)
└── Serves: Only units on floors 14-15
```

**Level 4 - Unit-Specific (Rare):**
```
"Personal Laundry Service"
├── Assigned to: Flat A-101 (resident's preferred vendor)
└── Serves: Only this specific flat
```

---

### Resident Filtering Logic

**How It Works:**

1. **Find Resident's Path:**
   - Resident in Flat A-101 has path: `1.2.4.6`
   - Ancestors: Society (1), Building A (2), Floor 1 (4), Flat A-101 (6)

2. **Match Vendor Assignments:**
   - Show vendors assigned to:
     - NULL (society-wide)
     - Node 1 (society root)
     - Node 2 (Building A)
     - Node 4 (Floor 1)
     - Node 6 (their flat)

3. **Path-Based Query:**
   ```
   Resident path '1.2.4.6' matches vendors assigned to any ancestor in that path
   ```

**Default View:**
- Resident sees vendors assigned to their hierarchy path
- Example: Flat A-101 resident sees vendors for Building A, Floor 1, or society-wide

**Override Option:**
- Resident can toggle "Show All Vendors" to see every approved vendor
- Useful for emergencies, personal preferences, or trying new services

---

### Vendor Order View Logic

**Vendor Dashboard:**
1. **Default:** Vendor sees all orders from units under their assigned nodes
2. **Filter Options:**
   - Filter by specific building/phase
   - Filter by service type
   - Filter by order status

**Example:**
- Vendor assigned to Building A
- Sees orders from: All floors and flats in Building A
- Can filter to: "Show only Floor 1 orders"

---

## 🔄 **Key Unique Features**

1. **Generic hierarchical model**: Flexible tree structure supports ANY society layout without schema changes
2. **Unlimited hierarchy depth**: Add as many levels as needed (Society → Building → Floor → Unit, or custom)
3. **Path-based queries**: Efficient ancestor/descendant lookups using PostgreSQL ltree extension
4. **Ultra-flexible vendor assignment**: Assign vendors to any hierarchy level (society, building, floor, or even individual unit)
5. **Hierarchical inheritance**: Vendor assigned to Building A automatically serves all floors and units within
6. **Smart vendor filtering**: Residents see vendors assigned to their path ancestors by default
7. **Override option**: Residents can view all vendors in society for emergencies or preferences
8. **No type constraints**: Works uniformly for apartments, layouts, mixed-use, and future society types
9. **Separate orders per category**: Can't mix laundry with car wash - each category is a separate order
10. **Mixed services within category**: One laundry order can have ironing + washing + dry cleaning items
11. **Independent workflow tracking**: Each service type follows its own completion steps
12. **Service-wise progress**: Ironing ready in 2 days while dry cleaning still processing (5 days)
13. **Manual payment confirmation**: Vendor self-reports payment receipt (CASH/UPI/CARD), auto-closes after 48h grace period
14. **Society subscription**: Societies pay platform monthly fee (₹5k-₹20k), vendors keep 100% of earnings
15. **Multi-category platform**: Built day 1 to support all categories, activate when ready
16. **Workflow flexibility**: Each service type can have unique completion steps
17. **Cross-category vendors**: One vendor can serve multiple categories with different workflows
18. **Zero rebuild needed**: Adding new categories/services OR new society structures = configuration, not development

---

## 📱 **Sample User Journeys**

### Current Implementation (All Categories Built, Only Laundry Active)

**Resident sees**:
```
Home Screen:
├─ 👔 Laundry Services [ACTIVE] → 50 providers
├─ 🚗 Vehicle Services [COMING SOON]
├─ 🏡 Home Services [COMING SOON]
└─ 💇 Personal Care [COMING SOON]

Taps Laundry → Creates order with mixed items
```

**Backend has**:
```sql
-- All category tables populated, only LAUNDRY set to is_live = true
parent_categories:
├─ LAUNDRY (is_live = true)  ✅ Active
├─ VEHICLE (is_live = false) 🔒 Ready but inactive
├─ HOME (is_live = false)    🔒 Ready but inactive
└─ PERSONAL (is_live = false) 🔒 Ready but inactive

-- When ready to launch Vehicle:
-- UPDATE parent_categories SET is_live = true WHERE category_key = 'VEHICLE'
-- Onboard vendors → Goes live immediately
```

---

### Detailed Journey Example

**Resident - Mixed Category Orders**:
```
Saturday 9 AM - Creates two separate orders:

ORDER #001 (Laundry Category):
├─ 5 Shirts (Ironing Only) @ ₹10 = ₹50
├─ 3 Pants (Washing + Ironing) @ ₹30 = ₹90
├─ 1 Suit (Dry Cleaning) @ ₹250 = ₹250
└─ Total: ₹390 | Pickup: Today 3 PM | Expected delivery: 5 days

ORDER #002 (Vehicle Category - when active):
├─ 1 Car Wash (Exterior + Interior) @ ₹300 = ₹300
├─ 1 Bike Wash @ ₹100 = ₹100
└─ Total: ₹400 | Service: Tomorrow 10 AM | Expected completion: 1 hour

Total spend across categories: ₹790
Two different vendors, two separate payments
```

**Tracking - Laundry Order #001**:
```
Day 1 (3 PM): Pickup complete
├─ Ironing workflow started: Pickup ✅ → Count ✅ → Iron ⏳
├─ Washing workflow started: Pickup ✅ → Count ✅ → Wash ⏳
└─ Dry Cleaning workflow started: Pickup ✅ → Count ✅ → Dry Clean ⏳

Day 2 (2 PM): Ironing complete
├─ Ironing: Pickup ✅ → Count ✅ → Iron ✅ → Ready ✅ → Deliver ⏳
├─ Washing: Pickup ✅ → Count ✅ → Wash ✅ → Iron ⏳
└─ Dry Cleaning: Pickup ✅ → Count ✅ → Dry Clean ⏳

Day 3 (4 PM): Washing complete
├─ Ironing: All steps ✅ (Waiting for full order)
├─ Washing: Pickup ✅ → Count ✅ → Wash ✅ → Iron ✅ → Ready ✅ → Deliver ⏳
└─ Dry Cleaning: Pickup ✅ → Count ✅ → Dry Clean ✅ → QC ⏳

Day 5 (5 PM): All complete - Delivery
├─ Ironing: All ✅
├─ Washing: All ✅
└─ Dry Cleaning: All ✅ → Single delivery of all items → Pay ₹390
```

**Tracking - Vehicle Order #002**:
```
Sunday 10:00 AM: Service starts
├─ Car Wash: Schedule ✅ → Arrive ✅ → Wash ⏳
└─ Bike Wash: Schedule ✅ → Waiting ⏳

Sunday 10:30 AM: Both in progress
├─ Car Wash: Wash ✅ → Vacuum ⏳
└─ Bike Wash: Wash ⏳

Sunday 10:50 AM: Completion
├─ Car Wash: All steps ✅ (Wash → Vacuum → Polish → Complete)
└─ Bike Wash: All steps ✅ (Wash → Dry → Complete)

Pay ₹400 → Both services complete
```

---

### Multi-Category Vendor Dashboard Example

**"QuickServe" offers Laundry + Vehicle services**

```
Today's Tasks:
├─ 👔 LAUNDRY ORDERS (8 orders, 15 service workflows)
│   ├─ Ironing workflows (5)
│   │   ├─ 2 at "Iron" step
│   │   └─ 3 at "Ready" step
│   ├─ Washing workflows (6)
│   │   ├─ 4 at "Wash" step
│   │   └─ 2 at "Iron" step
│   └─ Dry Cleaning workflows (4)
│       ├─ 2 at "Dry Clean" step
│       └─ 2 at "QC" step
│
└─ 🚗 VEHICLE ORDERS (4 orders, 6 service workflows)
    ├─ Car Wash workflows (4)
    │   ├─ 2 scheduled 10 AM
    │   └─ 2 scheduled 2 PM
    └─ Bike Wash workflows (2)
        └─ Both scheduled 11 AM

Revenue tracking by workflow:
├─ Ironing: ₹800 (40 items @ ₹20 avg)
├─ Washing: ₹1,200 (30 items @ ₹40 avg)
├─ Dry Cleaning: ₹2,000 (10 items @ ₹200 avg)
├─ Car Wash: ₹1,200 (4 cars @ ₹300)
└─ Bike Wash: ₹200 (2 bikes @ ₹100)

Total: ₹5,400 across 21 service workflows today
```

---

## 🎯 **System Architecture - Day 1 Build**

### Database includes ALL categories (even inactive ones)

```sql
-- Built from start:
parent_categories: LAUNDRY, VEHICLE, HOME, PERSONAL (all exist)

service_categories:
├─ Ironing, Washing, Dry Cleaning (LAUNDRY - active)
├─ Car Wash, Bike Wash, Detailing (VEHICLE - inactive)
├─ Gardening, Plumbing, Electrical (HOME - inactive)
└─ Barber, Salon, Spa (PERSONAL - inactive)

-- Workflow definitions exist for all:
service_workflows:
├─ Ironing: 5 steps defined
├─ Dry Cleaning: 6 steps defined
├─ Car Wash: 6 steps defined (ready to use)
├─ Gardening: 6 steps defined (ready to use)
└─ etc.
```

### Launch new category = 2 simple steps

1. **Activate category in database:**
   ```sql
   UPDATE parent_categories
   SET is_live = true
   WHERE category_key = 'VEHICLE';
   ```

2. **Onboard vendors** → Category goes live immediately

**No code changes. No schema changes. Just configuration.**

---

## 📊 **Order Structure**

### One Category = One Order

```javascript
// Resident creates separate orders per category
const laundryOrder = {
  order_id: "ORD001",
  category: "LAUNDRY",
  services: [
    { service_id: 1, items: [{ name: "Shirt", qty: 5, service: "Ironing" }] },
    { service_id: 3, items: [{ name: "Suit", qty: 1, service: "Dry Cleaning" }] }
  ],
  total: 350,
  payment_status: "pending"
}

const vehicleOrder = {
  order_id: "ORD002",
  category: "VEHICLE",
  services: [
    { service_id: 10, vehicle: "Car", service: "Car Wash" },
    { service_id: 11, vehicle: "Bike", service: "Bike Wash" }
  ],
  total: 400,
  payment_status: "pending"
}

// Two separate orders
// Two separate payments
// Two separate vendors (potentially)
// Two separate tracking workflows
```

---

## 🔧 **Service Workflow Configuration**

### Each service type has configurable workflow

```javascript
// Stored in database, configurable by admin
const serviceWorkflows = {
  IRONING: {
    steps: [
      { order: 1, name: "Pickup", required: true },
      { order: 2, name: "Count Items", required: true },
      { order: 3, name: "Iron", required: true },
      { order: 4, name: "Ready for Delivery", required: true },
      { order: 5, name: "Delivered", required: true }
    ],
    default_turnaround: 24, // hours
    pricing_model: "PER_ITEM"
  },

  DRY_CLEANING: {
    steps: [
      { order: 1, name: "Pickup", required: true },
      { order: 2, name: "Count Items", required: true },
      { order: 3, name: "Dry Clean", required: true },
      { order: 4, name: "Quality Check", required: true },
      { order: 5, name: "Ready for Delivery", required: true },
      { order: 6, name: "Delivered", required: true }
    ],
    default_turnaround: 120, // hours (5 days)
    pricing_model: "PER_ITEM"
  },

  CAR_WASH: {
    steps: [
      { order: 1, name: "Schedule", required: true },
      { order: 2, name: "Arrive on Site", required: true },
      { order: 3, name: "Exterior Wash", required: true },
      { order: 4, name: "Interior Vacuum", required: false },
      { order: 5, name: "Polish", required: false },
      { order: 6, name: "Complete", required: true }
    ],
    default_turnaround: 1, // hours
    pricing_model: "PER_SERVICE"
  },

  GARDENING: {
    steps: [
      { order: 1, name: "Schedule", required: true },
      { order: 2, name: "Arrive", required: true },
      { order: 3, name: "Trim Plants", required: false },
      { order: 4, name: "Mow Lawn", required: false },
      { order: 5, name: "Clean Up", required: true },
      { order: 6, name: "Complete", required: true }
    ],
    default_turnaround: 2, // hours
    pricing_model: "HOURLY"
  },

  PLUMBING: {
    steps: [
      { order: 1, name: "Schedule", required: true },
      { order: 2, name: "Arrive", required: true },
      { order: 3, name: "Diagnose Issue", required: true },
      { order: 4, name: "Fix", required: true },
      { order: 5, name: "Test", required: true },
      { order: 6, name: "Complete", required: true }
    ],
    default_turnaround: 3, // hours
    pricing_model: "PER_SERVICE"
  }
}

// Vendor updates progress:
updateServiceProgress(orderId, serviceId, currentStep);

// Resident sees real-time progress per service type
```

---

## 📈 **Platform Evolution - Zero Rebuild**

### Timeline

**Month 1-3**: Launch with Laundry only
- Database: All categories exist but only LAUNDRY is_live = true
- UI: Shows "Coming Soon" for other categories
- Vendors: Can only register for Laundry services
- Orders: Only laundry orders possible

**Month 4**: Activate Vehicle Services
- Action: `SET is_live = true` for VEHICLE category
- Onboard: Car wash and bike wash vendors
- Launch: Vehicle services go live
- Development: 0 hours (just configuration)

**Month 6**: Activate Home Services
- Action: `SET is_live = true` for HOME category
- Onboard: Gardening, plumbing, electrical vendors
- Launch: Home services go live
- Development: 0 hours (just configuration)

**Month 9**: Activate Personal Care
- Action: `SET is_live = true` for PERSONAL category
- Onboard: Barber, salon, spa vendors
- Launch: Personal care services go live
- Development: 0 hours (just configuration)

### Technical Effort Per New Category

```
Development time: 0 hours
Schema changes: 0
Code changes: 0
Testing: Functional testing only
Deployment: Configuration update

Process:
1. Flip is_live flag in database
2. Onboard vendors for that category
3. Market to residents
4. Monitor and optimize
```

---

## 💡 **Business Model**

### Revenue

**Society Subscription**:
- Starter (100-300 flats): ₹5,000/month
- Growth (301-600 flats): ₹10,000/month
- Enterprise (601+ flats): ₹20,000/month

**Covers ALL service categories** (Laundry, Vehicle, Home, Personal)

### Vendor Earnings

**Vendors keep 100% of order value**
- Resident pays vendor directly (UPI/Cash)
- No commission deducted
- No transaction fees
- Immediate payment settlement

### Win-Win-Win

**Residents**: All home services in one app
**Vendors**: 100% earnings + access to society customers
**Societies**: Fixed monthly fee + organized vendor management
**Platform**: Predictable recurring revenue

---

## 🎬 **End State Vision**

**Comprehensive Society-Based Home Services Marketplace**

One platform where:
- **Residents** get all home services (laundry, vehicle care, home maintenance, personal care)
- **Vendors** serve multiple categories with unified workflow management
- **Societies** manage all vendors and services with single subscription
- **Platform** scales infinitely by activating new categories

**Zero technical debt. Built right from Day 1.**

---

**End of Document**
