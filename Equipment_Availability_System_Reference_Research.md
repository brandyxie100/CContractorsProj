# Equipment Availability & Rental Management System — Reference Research Document

**Project**: Clements Contractors Project
**Sub-project**: Equipment Fleet Availability & Rental Management Module  
**Target Customer**: Clements Contractors (Earthworks & Contracting)  
**Document Date**: 2026-08-08  
**Status**: v1.1 — For Product Planning & Design Reference

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Key Availability Features Comparison](#2-key-availability-features-comparison)
3. [Target Customer Profile](#3-target-customer-profile)
4. [Primary Reference: Clue](#4-primary-reference-clue)
5. [UI Reference: EZRentOut Availability Calendar](#5-ui-reference-ezrentout-availability-calendar)
6. [Differentiator Reference: RentalMan IntelliSource](#6-differentiator-reference-rentalman-intellisource)
7. [Competitive Feature Matrix](#7-competitive-feature-matrix)
8. [Key Design Patterns & UX Insights](#8-key-design-patterns--ux-insights)
9. [Recommended Feature Roadmap](#9-recommended-feature-roadmap)
10. [Technical Architecture Considerations](#10-technical-architecture-considerations)
11. [Next Steps](#11-next-steps)
12. [Appendix: Reference Links](#appendix-reference-links)

---

## 1. Executive Summary

This document compiles research on leading equipment rental and fleet management software systems, with the goal of informing the design and development of an equipment availability tracking module for our SaaS platform. The target customer is **Clements Contractors**, an earthworks and contracting company with their own fleet of trucks, diggers, and quarry equipment.

Three reference systems have been selected for different purposes:

| Reference System | Role | Why Selected |
|---|---|---|
| **Clue** | Primary reference | The only system built *for contractors* (not rental companies) — matches our target customer profile exactly |
| **EZRentOut** | UI/UX reference | Clean, well-documented availability calendar — a pattern users already understand |
| **RentalMan (Wynne Systems)** | Differentiator reference | IntelliSource feature — smart "next available" search that can set our product apart |

The key insight: most equipment rental software is built for *rental companies* who rent equipment out to others. Clements Contractors is a *contractor* who owns their fleet and deploys it across their own job sites. This distinction matters for workflow, pricing model, and feature priorities.

---

## 2. Key Availability Features Comparison

This section provides a focused, side-by-side breakdown of how each platform handles **equipment availability** — the core capability that matters most for fleet management.

### 2.1 Availability Tracking Capabilities

| Capability | Clue | EZRentOut | RentalMan (IntelliSource) |
|---|---|---|---|
| **Real-time availability status** | ✅ Per-asset status with live updates | ✅ Per-asset and per-category | ✅ Enterprise-grade, cross-branch |
| **Next-available date/time** | ✅ Shows when equipment becomes free | ✅ Visible in calendar and list views | ✅ Smart search returns exact availability |
| **Cross-location availability** | ✅ Across all sites and depots | ✅ Multi-location inventory | ✅ Cross-branch with IntelliSource |
| **Maintenance-aware availability** | ✅ Auto-blocks when service is due | ✅ Maintenance status blocks booking | ✅ Service windows factored in |
| **Quantity-level availability** | ⚠️ Limited (asset-centric) | ✅ "5 units available" counters | ✅ Both quantity and individual |
| **Individual asset-level tracking** | ✅ Serialized, per-machine | ✅ Serialized assets | ✅ Full asset lifecycle tracking |
| **GPS / telematics integration** | ✅ Samsara, John Deere, Caterpillar | ⚠️ Limited / third-party only | ✅ Telematics feeds availability |
| **Utilization-based availability** | ✅ Tied to usage hours | ⚠️ Basic | ✅ Advanced (meter-based) |

### 2.2 Availability Display & UI

| Feature | Clue | EZRentOut | RentalMan |
|---|---|---|---|
| **Calendar timeline view** | ✅ Gantt-style timeline | ✅ Industry-leading calendar UI | ✅ Standard rental calendar |
| **Color-coded statuses** | ✅ Status badges + colors | ✅ 5+ status colors | ✅ Status color coding |
| **Drag-and-drop scheduling** | ✅ Drag to assign/reschedule | ✅ Drag to create reservations | ✅ Drag-and-drop dispatch |
| **Day / week / month views** | ✅ All three | ✅ All three + custom ranges | ✅ All three |
| **List view with filters** | ✅ Sortable, filterable asset list | ✅ Advanced filtering | ✅ Enterprise-grade filtering |
| **Map view** | ⚠️ Via telematics integration | ❌ Not native | ⚠️ Limited |
| **Mobile app** | ✅ Native iOS / Android | ✅ Mobile app + responsive | ✅ MobileLink field app |
| **Dashboard widgets** | ✅ Customizable dashboard | ✅ Dashboard with KPIs | ✅ Enterprise dashboards |

### 2.3 Smart Availability Features (Differentiators)

| Feature | Clue | EZRentOut | RentalMan |
|---|---|---|---|
| **"Find me a machine" smart search** | ⚠️ Basic filtering | ❌ No native smart search | ✅ **IntelliSource** — flagship feature |
| **Ranked results (by availability/proximity)** | ⚠️ Partial | ❌ | ✅ Ranked by best match |
| **Alternative equipment suggestions** | ❌ | ❌ | ✅ Suggests similar equipment |
| **Nearest asset finder** | ⚠️ Via GPS integration | ❌ | ✅ Built into IntelliSource |
| **Transport cost / distance estimate** | ❌ | ❌ | ⚠️ Partial (logistics module) |
| **Conflict detection (double-booking)** | ✅ Real-time conflict alerts | ✅ Automatic conflict prevention | ✅ Enterprise-grade conflict engine |
| **One-click reservation from search** | ⚠️ | ✅ Quick-book from calendar | ✅ Reserve directly from search |
| **What-if scenario planning** | ❌ | ❌ | ⚠️ Advanced scheduling module |

### 2.4 Availability-Related Workflows

| Workflow | Clue | EZRentOut | RentalMan |
|---|---|---|---|
| **Equipment check-in / check-out** | ✅ Field app + scan | ✅ QR / barcode scanning | ✅ MobileLink scanning |
| **Reservation / booking management** | ✅ Project-based assignments | ✅ Full rental reservation system | ✅ Enterprise contract management |
| **Inter-location transfers** | ✅ Transfer tracking | ✅ Branch transfers | ✅ Cross-branch sourcing |
| **Maintenance blocks availability** | ✅ Auto-removes from availability | ✅ Service orders block booking | ✅ PM scheduling auto-blocks |
| **Utilization reporting** | ✅ Detailed utilization analytics | ✅ Utilization dashboards | ✅ Enterprise fleet analytics |
| **Bundles / kits management** | ❌ | ✅ Bundle multiple items as one | ⚠️ Limited (package rentals) |
| **Availability alerts / notifications** | ✅ Email + in-app alerts | ✅ Automated reminders | ✅ Enterprise alerting engine |

### 2.5 Summary: Who Does What Best

| Category | Winner | Why |
|---|---|---|
| **Best for contractors** | **Clue** | Built for contractors, not rental companies. Project-centric, telematics-first. |
| **Best calendar UI** | **EZRentOut** | Cleanest, most intuitive availability calendar. Industry-standard pattern. |
| **Best smart availability** | **RentalMan** | IntelliSource is the gold standard for "find me the right machine" search. |
| **Best mobile experience** | **Clue** | Field app built for on-site crews and operators. |
| **Best multi-location** | **RentalMan** | Enterprise-grade cross-branch availability and sourcing. |

---

## 3. Target Customer Profile

### Clements Contractors — Company Overview

- **Industry**: Earthworks & Contracting
- **Years in Business**: 28+ years
- **Core Services**:
  - Driveway surfacing
  - Concrete works
  - Drainage
  - Earthworks & excavation
  - Cartage/deliveries (large & small trucks)
  - Land clearing & development
  - Landscaping & supplies
  - Quarry materials
  - Retaining walls
  - Roading & subdivisions
  - Road maintenance
  - Traffic management
- **Key Assets**: Own quarries, fleet of trucks, superb equipment, experienced crew
- **Business Model**: Contractor that owns and operates its own equipment across multiple job sites

### Pain Points to Solve

Based on their business profile, the key equipment management challenges are likely:

1. **Fleet visibility across multiple job sites** — Knowing where each truck/digger is and when it will be free for the next job.
2. **Scheduling conflicts** — Double-booking equipment or having crews wait because a machine is still on another site.
3. **Maintenance coordination** — Ensuring equipment is serviced on time without disrupting project schedules.
4. **Utilization tracking** — Understanding which equipment is underutilized or overworked to inform buying/renting decisions.
5. **Dispatch efficiency** — Quickly finding the nearest available machine for an urgent job.

---

## 4. Primary Reference: Clue

### 4.1 Overview

**Website**: getclue.com  
**Tagline**: "The best software for equipment management"  
**Target Market**: Construction contractors, heavy civil, earthworks  
**Pricing**: Contact vendor (enterprise)

Clue is a construction equipment management platform designed specifically for contractors — not rental companies. It integrates fleet tracking, maintenance, dispatch, rentals, and field workflows into a single system.

### 4.2 Why Clue Is Our Primary Reference

Clue is the only major system in this space built from the ground up for **contractors who own their own fleet and deploy it across their own job sites**. Every other major player (EZRentOut, Point of Rental, RentalMan) is built for *rental companies* that rent equipment out to third-party customers.

This distinction is critical because:
- **No invoicing/billing for rentals** — Contractors don't charge themselves rent. The system focuses on utilization, cost tracking, and dispatch, not rental revenue.
- **Project-centric, not customer-centric** — Equipment is allocated to *projects*, not rented to *customers*.
- **Internal workflow focus** — The value is in operational efficiency, not revenue generation from rentals.
- **Telematics-first** — Deep integration with GPS and engine hour tracking, which is essential for heavy equipment.

### 4.3 Core Features

#### 4.3.1 Unified Fleet Visibility
- Real-time tracking of all equipment (owned, rented, managed)
- Telematics integrations with **Samsara, John Deere, and Caterpillar**
- Shows location, status, and availability across all projects
- Single dashboard view of the entire fleet

#### 4.3.2 Intelligent Dispatch Management
- **Drag-and-drop scheduling** — Assign equipment to jobs with a visual interface
- Matches equipment requests with available assets based on **proximity and utilization**
- Mobile apps enable instant field updates and reassignments
- Field teams can request equipment reassignments from the job site

#### 4.3.3 Automated Preventive Maintenance
- Scheduling automatically accounts for **maintenance windows** based on actual usage hours
- Creates work orders when service is due
- **Blocks assignments** until maintenance is completed
- Full service history tracking per asset

#### 4.3.4 Advanced Utilization Analytics
- Comprehensive reporting on:
  - Idle time
  - Productive hours
  - Utilization rates
  - Cost per project
- Identifies underutilized equipment
- Informs fleet composition decisions (buy vs. rent)

#### 4.3.5 Mobile Field App
- On-site teams can:
  - Request equipment reassignments
  - Mark equipment for return
  - Report issues
  - View schedules
- Real-time sync with the office dashboard

### 4.4 UI/UX Patterns to Reference

#### Dashboard Layout
- **Left sidebar navigation**: Home, Performance Dashboard, Track Everything, Resource Planner (Planning, Dispatch), Project Productivity, Reports, Maintenance (Fault Codes, Inspection Issues, Work Orders, Preventive Maintenance), Directory (Assets, People, Projects, Geofence)
- **Top bar**: Quick search, notifications, user profile
- **Main content area**: Data tables, charts, and visualizations

#### Availability Status Indicators
- Color-coded status badges (e.g., "AVAILABLE" in blue/green)
- Asset ID and model number prominently displayed
- Current job assignment shown inline
- Expected return date visible at a glance

#### Maintenance Dashboard
- Three-column status overview: Overdue / Upcoming / Good Standing
- Asset list with PM plan, PM task, PM interval, current usage
- Progress bars showing hours until next service
- Color-coded urgency indicators

### 4.5 Key Takeaways for Our Product

1. **Project-centric data model** — Equipment is assigned to projects/jobs, not rented to customers.
2. **Telematics integration is table stakes** — For heavy equipment, GPS and hour meter data are essential.
3. **Maintenance drives availability** — The system must automatically factor maintenance windows into availability calculations.
4. **Mobile is not optional** — Field crews need to interact with the system from job sites.
5. **Utilization analytics is the value driver** — For contractors, the ROI comes from better utilization, not rental revenue.

---

## 5. UI Reference: EZRentOut Availability Calendar

### 5.1 Overview

**Website**: ezo.com/ezrentout  
**Target Market**: Equipment rental companies (IT, AV, tool, construction)
**Pricing**: From $89/month (Essential plan)
**G2 Rating**: Strong, widely used across industries

EZRentOut is a cloud-based equipment rental management platform. While it's built for rental companies (not contractors), its **availability calendar UI** is one of the cleanest, most intuitive implementations in the industry — and a pattern that users across the sector already understand.

### 5.2 Why EZRentOut for UI Reference

- **Mature, well-documented UI** — Thousands of users are familiar with the pattern
- **Clean visual design** — The availability calendar is easy to read at a glance
- **Proven UX pattern** — The calendar-based availability view is the industry standard
- **Great for onboarding** — Users who have used any rental software will feel familiar

### 5.3 Availability Calendar — Detailed Breakdown

#### 5.3.1 Core Calendar Features

| Feature | Description |
|---|---|
| **Visual timeline view** | Horizontal timeline showing each asset's availability over days/weeks |
| **Color-coded statuses** | Different colors for available, booked, maintenance, transit |
| **Consolidated view** | See all equipment availability in one calendar |
| **Quantity tracking** | Shows exact number of units available at any time |
| **Conflict detection** | Automatically flags double-bookings |
| **Drag-to-create** | Create reservations by dragging on the calendar |
| **Filtering** | Filter by equipment type, location, status |

#### 5.3.2 Status States (Color Coding)

The calendar typically uses these status states (we should adopt a similar pattern):

| Status | Color (typical) | Meaning |
|---|---|---|
| Available | Green / White | Ready to deploy |
| On Rent / Assigned | Blue | Currently on a job site |
| Maintenance | Orange / Red | Out of service for maintenance |
| Transit | Yellow | Being moved between sites |
| Reserved | Light blue / Purple | Booked for an upcoming job |

#### 5.3.3 Calendar Views

- **Day view** — Hour-by-hour breakdown (for short-term rentals/dispatches)
- **Week view** — 7-day overview (most common for planning)
- **Month view** — High-level long-range planning
- **Gantt-style** — Horizontal bars showing duration of each assignment

### 5.4 Additional EZRentOut Features Worth Noting

#### 5.4.1 Equipment Tracking
- Records each piece of equipment's location, usage history, and rental status
- Helps teams quickly identify where equipment is deployed and when it will become available
- Prevents lost equipment, reduces idle time, improves fleet utilization

#### 5.4.2 Bundles & Kits Management
- Construction equipment is often rented/deployed in groups
- Example: A concrete job might require mixers, power tools, safety gear, and attachments
- Bundles let you manage multiple related items as a single unit

#### 5.4.3 QR / Barcode Scanning
- Label equipment with QR codes
- Crews scan to check in/check out equipment
- Streamlines the equipment handoff process
- Creates audit trail of who had what equipment and when

#### 5.4.4 Webstore / Customer Portal
- Customizable online booking interface
- Customers can browse catalog and check availability
- Less relevant for contractors (internal use only), but the self-service pattern is valuable

### 5.5 Key UI/UX Takeaways for Our Product

1. **Calendar is the central view** — The availability calendar should be the first thing users see when they log in. It's the most information-dense and useful view.
2. **Color coding is essential** — Users should be able to tell equipment status at a glance, without reading text.
3. **Drag-and-drop assignment** — Let users create assignments by dragging on the calendar. It's intuitive and fast.
4. **Filter by type/location** — With a large fleet, users need to filter down to what they care about (e.g., "show me all excavators").
5. **Quantity vs. individual assets** — Support both: category-level availability ("3 excavators available") and individual asset-level tracking ("Excavator #7 is available").

---

## 6. Differentiator Reference: RentalMan IntelliSource

### 6.1 Overview

**Product**: RentalMan by Wynne Systems  
**Website**: wynnesystems.com  
**Target Market**: Large construction equipment rental companies, heavy equipment fleets  
**Tier**: Enterprise  
**Pricing**: Contact vendor

RentalMan is an enterprise-grade rental ERP platform. While it's overkill for a small-to-mid contractor, it has one standout feature that we should study closely: **IntelliSource**.

### 6.2 What Is IntelliSource?

IntelliSource is RentalMan's **smart equipment sourcing and availability engine**. It answers the question: *"Where is the nearest available excavator, and when will it be free?"*

Instead of manually scrolling through calendars or checking each location, a dispatcher can ask the system to find available equipment across the entire fleet, across all locations.

### 6.3 How IntelliSource Works

#### 6.3.1 Core Capabilities

| Capability | Description |
|---|---|
| **Cross-location search** | Finds available equipment across all branches/depots, not just one location |
| **Scheduling conflict detection** | Automatically flags when equipment is already booked |
| **Alternative suggestions** | If the exact asset isn't available, suggests similar equipment that is |
| **Nearest asset finder** | Identifies the closest available unit to minimize transport time/cost |
| **Availability timeline** | Shows *when* equipment will become available (not just if it's available now) |
| **Maintenance-aware** | Factors in upcoming service when calculating availability |

#### 6.3.2 Example Workflow

A project manager needs a 20-ton excavator for a job starting next Monday:

1. **PM opens IntelliSource** and enters: "20-ton excavator, needed from Monday 8am to Wednesday 5pm, job site at Location A"
2. **System searches entire fleet** across all yards/depots
3. **Results show**:
   - 2 excavators available immediately at Depot B (15 min from site)
   - 1 excavator available Tuesday morning at Depot C (30 min from site)
   - 1 excavator currently on rent at Job X, due back Monday afternoon at Depot A (5 min from site)
4. **PM selects the best option** and reserves it with one click

### 6.4 Why This Is a Differentiator

Most equipment management systems show you a **calendar** — you look at it and figure out availability yourself. IntelliSource flips the model: you tell the system *what you need*, and it tells you *what's available, where, and when*.

This is a meaningful upgrade for several reasons:
- **Faster decision-making** — Dispatchers don't need to scan dozens of assets manually
- **Fewer mistakes** — The system doesn't miss an available asset in a different location
- **Better utilization** — Equipment gets moved to where it's needed more efficiently
- **Scales with fleet size** — The value increases as the fleet grows (manual checking gets harder)

### 6.5 How We Can Adapt This for Our Product

For a contractor-focused product (not a rental company), we'd adapt IntelliSource into a **"Smart Dispatch Finder"**:

#### 6.5.1 Smart Dispatch Finder — Proposed Features

| Feature | Description |
|---|---|
| **"Find me a machine" search** | User inputs: equipment type, date range, job site location |
| **Availability ranking** | Results ranked by: earliest availability, closest location, lowest cost |
| **Next-available date** | For each asset, shows the exact date/time it becomes free |
| **Alternative suggestions** | If exact match unavailable, suggests similar equipment (e.g., "No 20t excavator, but 25t is available") |
| **Transport cost estimate** | Shows estimated transport time/cost from current location to job site |
| **One-click assignment** | Reserve the selected asset directly from search results |
| **Conflict prevention** | Automatically checks for overlapping assignments |

#### 6.5.2 Example User Flow

```
User: "I need an excavator for the Smith Street job, starting tomorrow."

System shows:
┌─────────────────────────────────────────────────────────────┐
│  🚜 Excavator #7 (20t) — AVAILABLE NOW                      │
│  Location: Depot (12 min from Smith St)                     │
│  Next assignment: None scheduled                            │
│  [Assign to Job]                                             │
├─────────────────────────────────────────────────────────────┤
│  🚜 Excavator #3 (25t) — AVAILABLE TOMORROW 2PM             │
│  Location: Jones Road job (8 min from Smith St)             │
│  Next assignment: None after tomorrow 2PM                   │
│  [Assign to Job]                                             │
├─────────────────────────────────────────────────────────────┤
│  🚜 Excavator #5 (15t) — AVAILABLE NOW                      │
│  Location: Depot (12 min from Smith St)                     │
│  ⚠️ Smaller than requested (15t vs 20t)                     │
│  [Assign to Job]                                             │
└─────────────────────────────────────────────────────────────┘
```

### 6.6 Technical Implementation Notes

To build a Smart Dispatch Finder, we'd need:

1. **Asset database** with type, specifications, location, and status
2. **Assignment/booking database** with start/end times and job site locations
3. **Maintenance schedule database** with service windows
4. **Geolocation data** for distance calculations (either GPS telematics or manual location entries)
5. **Search algorithm** that:
   - Filters by equipment type/specs
   - Checks availability for the date range
   - Calculates distance from current location to job site
   - Ranks results by relevance (availability + proximity + suitability)
   - Suggests alternatives if exact match not found

---

## 7. Competitive Feature Matrix

| Feature | Clue | EZRentOut | Rentman | Point of Rental | RentalMan | Our Product (Target) |
|---|---|---|---|---|---|---|
| **Built for contractors** | ✅ | ❌ | ❌ | ❌ | ❌ | ✅ |
| **Availability calendar** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Next-available date** | ✅ | ✅ | ✅ (minute-level) | ✅ | ✅ (IntelliSource) | ✅ |
| **Smart availability search** | ⚠️ Basic | ❌ | ❌ | ⚠️ Basic | ✅ (IntelliSource) | ✅ (Smart Dispatch Finder) |
| **GPS/telematics integration** | ✅ | ⚠️ Limited | ❌ | ✅ | ✅ | ✅ (phase 2) |
| **Maintenance scheduling** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Maintenance blocks availability** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Multi-site/location support** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Mobile field app** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ (phase 2) |
| **QR/barcode scanning** | ⚠️ | ✅ | ✅ | ✅ | ✅ | ⚠️ (phase 2) |
| **Utilization analytics** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Rental billing/invoicing** | ❌ | ✅ | ✅ | ✅ | ✅ | ❌ |
| **Customer portal** | ❌ | ✅ | ✅ | ✅ | ✅ | ❌ |
| **Project-centric data model** | ✅ | ❌ | ❌ | ❌ | ❌ | ✅ |
| **Public/transparent pricing** | ❌ | ✅ | ✅ | ❌ | ❌ | ✅ |

**Legend**: ✅ = Full support, ⚠️ = Partial/limited, ❌ = Not available

---

## 8. Key Design Patterns & UX Insights

### 8.1 Information Architecture

The system should be organized around these core modules:

1. **Dashboard** — High-level fleet overview, key metrics, alerts
2. **Availability Calendar** — Central view of all equipment schedules
3. **Dispatch / Assignments** — Manage equipment assignments to jobs
4. **Assets / Fleet** — Equipment inventory, specs, details
5. **Maintenance** — Service schedules, work orders, history
6. **Projects / Jobs** — Job sites, project timelines
7. **Reports / Analytics** — Utilization, costs, productivity
8. **Settings** — Users, locations, custom fields

### 8.2 Availability Visualization Patterns

#### Pattern 1: Calendar Timeline View (Primary)
- Horizontal timeline with assets listed vertically
- Color-coded bars showing assignments/maintenance
- Drag to create/resize assignments
- Zoom in/out for day/week/month views
- **Best for**: Detailed scheduling, dispatch planning

#### Pattern 2: List View with Status Badges
- Table of all assets with current status, location, next available date
- Sortable by status, type, location, availability date
- Filterable by equipment type, site, status
- **Best for**: Quick lookups, finding specific equipment

#### Pattern 3: Card / Grid View
- Each asset shown as a card with photo, ID, status badge, location
- Color-coded border or status indicator
- **Best for**: Visual overview, mobile views, less technical users

#### Pattern 4: Map View (Advanced)
- All assets plotted on a map
- Color-coded pins showing status
- Click to see details and availability
- **Best for**: Dispatchers optimizing for location, large geographic areas

### 8.3 Status System Design

We recommend a 5-state availability model:

| Status | Color | Meaning | Can be assigned? |
|---|---|---|---|
| Available | Green | Ready to deploy, at depot | ✅ Yes |
| Assigned / On Job | Blue | Currently deployed to a job site | ❌ No (but can schedule for after) |
| In Transit | Yellow | Being moved between locations | ❌ No (temporary) |
| Maintenance | Orange | Out of service for maintenance/repair | ❌ No |
| Reserved | Purple | Booked for an upcoming job | ⚠️ Only if booking doesn't conflict |

### 8.4 "Next Available" Display Pattern

For each asset, always show:
- **Current status** (with color)
- **Current location** (depot or job site name)
- **Next available date/time** (e.g., "Available: Aug 12, 5:00 PM")
- **Next assignment** (if booked: "Next: Jones Road Subdivision, Aug 13–18")

This gives dispatchers the key information at a glance without needing to open a detail view.

---

## 9. Recommended Feature Roadmap

### Phase 1: MVP (Core Availability Tracking)
**Goal**: Replace spreadsheets and whiteboards with a digital availability system

| Feature | Priority | Description |
|---|---|---|
| Asset registry | P0 | Add all equipment with type, model, ID, specs, photos |
| Availability calendar | P0 | Visual timeline showing all equipment assignments |
| Create/edit assignments | P0 | Assign equipment to jobs with start/end dates |
| Project/job management | P0 | Create job sites/projects to assign equipment to |
| Status tracking | P0 | Available / Assigned / Maintenance statuses |
| Next-available date | P0 | Show when each piece of equipment becomes free |
| Basic filtering | P0 | Filter by equipment type, status, location |
| User accounts & roles | P0 | Admin / Manager / Viewer roles |

### Phase 2: Smart Features & Operations
**Goal**: Add intelligence and operational efficiency

| Feature | Priority | Description |
|---|---|---|
| **Smart Dispatch Finder** | P1 | "Find me an excavator for tomorrow" — search across fleet |
| Alternative suggestions | P1 | Suggest similar equipment if exact match unavailable |
| Maintenance scheduling | P1 | Service reminders, work orders, maintenance blocks availability |
| Utilization reports | P1 | Hours used, idle time, utilization rate per asset |
| Mobile-responsive web app | P1 | Field crews can view schedules on phones |
| Location/distance tracking | P2 | Track equipment location, calculate transport times |
| Calendar drag-and-drop | P2 | Create/resize assignments by dragging on calendar |
| Notifications & alerts | P2 | Email/SMS for upcoming assignments, maintenance due |

### Phase 3: Advanced & Integrations
**Goal**: Full-featured fleet management platform

| Feature | Priority | Description |
|---|---|---|
| Telematics integration | P2 | GPS tracking, engine hour meters, fault codes |
| QR code check-in/out | P2 | Scan equipment to assign/return |
| Cost tracking | P2 | Track cost per hour, per job, per project |
| Equipment bundles/kits | P3 | Group related equipment (e.g., excavator + attachments) |
| API & integrations | P3 | Connect to accounting, project management, HR systems |
| Custom reporting | P3 | Build custom reports and dashboards |
| Multi-company / multi-branch | P3 | Support for companies with multiple divisions |

---

## 10. Technical Architecture Considerations

### 10.1 Data Model (Core Entities)

```
Asset (Equipment)
├── id
├── asset_id (human-readable: "EXC-007")
├── name / model
├── type (excavator, truck, roller, etc.)
├── category (heavy equipment, trucks, tools, etc.)
├── specs (weight, capacity, attachments)
├── status (available, assigned, maintenance, transit)
├── current_location_id
├── purchase_date
├── purchase_cost
├── hourly_rate (internal cost)
├── photo_url
├── qr_code
├── created_at
└── updated_at

Project / Job
├── id
├── name
├── address / location
├── start_date
├── end_date
├── status (planning, active, completed)
├── project_manager_id
└── description

Assignment
├── id
├── asset_id
├── project_id
├── start_time
├── end_time
├── status (scheduled, active, completed, cancelled)
├── assigned_by (user_id)
├── notes
└── created_at

Maintenance Record
├── id
├── asset_id
├── type (preventive, repair, inspection)
├── description
├── scheduled_date
├── completed_date
├── status (scheduled, in_progress, completed, overdue)
├── cost
├── performed_by
└── notes

Location
├── id
├── name
├── type (depot, job_site, workshop)
├── address
├── lat / lng
└── description

User
├── id
├── name
├── email
├── role (admin, manager, operator, viewer)
├── phone
└── avatar_url
```

### 10.2 Availability Calculation Logic

To determine if an asset is available for a given date range:

```
is_available(asset_id, start_time, end_time):
  1. Check for overlapping assignments:
     - Any assignment where:
       assignment.start_time < end_time AND assignment.end_time > start_time
     - If overlap exists → NOT AVAILABLE
     
  2. Check for overlapping maintenance:
     - Any maintenance record where:
       maintenance.scheduled_date < end_time AND maintenance.completed_date > start_time
       (or maintenance is in_progress)
     - If overlap exists → NOT AVAILABLE (MAINTENANCE)
     
  3. Check asset status:
     - If status is "retired" or "sold" → NOT AVAILABLE
     
  4. If no conflicts → AVAILABLE
```

### 10.3 "Next Available" Calculation

```
next_available_date(asset_id):
  1. Get current time
  2. Find the assignment that ends latest but starts before now
     (the current assignment)
  3. If no current assignment:
     - Check if in maintenance → return maintenance completion date
     - Otherwise → return "now" (available immediately)
  4. If there is a current assignment:
     - Look for any maintenance scheduled right after the assignment ends
     - Return the later of: assignment end_date OR maintenance end_date
  5. Also check for any future assignments that start soon
     - Return the gap between current assignment end and next assignment start
```

### 10.4 Smart Dispatch Finder Algorithm

```
find_available_equipment(equipment_type, start_time, end_time, job_location):
  1. Filter all assets by equipment_type
  2. For each asset, check availability for the date range
  3. For available assets:
     a. Calculate distance from current location to job_location
     b. Calculate transport time/cost
     c. Calculate suitability score (how well specs match the request)
  4. For unavailable assets:
     a. Find next_available_date
     b. Calculate how close it is to the requested start_time
  5. Sort results by:
     - Available now + closest location (highest priority)
     - Available soon + close location
     - Available later
     - Alternative equipment (different type but similar capability)
  6. Return ranked list with availability dates and transport info
```

---

## 11. Next Steps

### Immediate (This Week)
- [ ] Review this document with the product team
- [ ] Validate assumptions about Clements Contractors' needs
- [ ] Schedule demo/trial of Clue to see it in action
- [ ] Create wireframes for the availability calendar view
- [ ] Define MVP feature set more precisely

### Short Term (Next 2–4 Weeks)
- [ ] Design the core data model and database schema
- [ ] Build low-fidelity prototypes of key screens
- [ ] Conduct user interviews with earthworks contractors
- [ ] Define the UI design system (colors, components, typography)
- [ ] Start technical architecture planning

### Medium Term (1–3 Months)
- [ ] Build MVP (Phase 1 features)
- [ ] Internal testing and iteration
- [ ] Pilot with Clements Contractors
- [ ] Gather feedback and refine
- [ ] Plan Phase 2 features

### Open Questions
1. Does Clements Contractors currently use any software for equipment tracking? (If so, what and what do they dislike about it?)
2. How many pieces of equipment do they have? (Fleet size affects UI/UX decisions)
3. How many job sites do they typically run simultaneously?
4. Do they have telematics/GPS on their equipment already?
5. Who are the primary users? (Dispatchers, project managers, equipment managers, owners?)
6. What's their current process for scheduling equipment? (Spreadsheets? Whiteboards? Phone calls?)

---

## Appendix: Reference Links

- **Clue**: https://www.getclue.com/
- **EZRentOut**: https://ezo.com/ezrentout/
- **RentalMan (Wynne Systems)**: https://wynnesystems.com/
- **Rentman**: https://rentman.io/
- **Point of Rental**: https://www.pointofrental.com/
- **Reservety Construction Equipment Rental Guide**: https://reservety.com/guides/construction-equipment/construction-equipment-rental-software.html

---

*Document prepared for internal product planning purposes. All product names and trademarks belong to their respective owners.*
