const API = `${window.location.protocol}//${window.location.hostname}:8082/api`;
let selectedRoute;

const $ = (id) => document.getElementById(id);
function toast(message) { const el = $("toast"); el.textContent = message; el.classList.add("show"); setTimeout(() => el.classList.remove("show"), 3000); }
function renderRoutes(items) {
  $("results").innerHTML = items.length ? items.map((route) => `
    <article class="trip-card">
      <span class="transport-icon">${route.vehicle_type === "BUS" ? "▣" : "↝"}</span>
      <div><div class="trip-title">${route.name} <span class="muted">${route.vehicle_number}</span></div>
      <div class="trip-meta">${route.source} → ${route.destination}</div><div class="trip-time">${route.departure_time}</div></div>
      <button class="secondary-button" data-route="${route.id}">View seats</button>
    </article>`).join("") : `<div class="booking-result">No routes found. Try a different route or transport type.</div>`;
  $("result-count").textContent = `${items.length} route${items.length === 1 ? "" : "s"}`;
  $("results").querySelectorAll("[data-route]").forEach((button) => button.addEventListener("click", () => openBooking(Number(button.dataset.route))));
}
async function openBooking(id) {
  selectedRoute = window.searchResults.find((route) => route.id === id);
  $("dialog-title").textContent = `${selectedRoute.name} · ${selectedRoute.source} → ${selectedRoute.destination}`;
  const response = await fetch(`${API}/seats?vehicleId=${id}`);
  if (!response.ok) return toast("Unable to load seats.");
  const data = await response.json();
  $("seat-grid").innerHTML = data.seats.map((seat) => `<button type="button" class="seat ${seat.status !== "AVAILABLE" ? "booked" : ""}" data-seat="${seat.id}" ${seat.status !== "AVAILABLE" ? "disabled" : ""}>${seat.seat_number}</button>`).join("");
  $("seat-grid").querySelectorAll("[data-seat]").forEach((button) => button.addEventListener("click", () => {
    $("seat-grid").querySelectorAll(".selected").forEach((item) => item.classList.remove("selected")); button.classList.add("selected"); $("selected-seat").value = button.dataset.seat;
  }));
  $("booking-dialog").showModal();
}
$("search-form").addEventListener("submit", (event) => {
  event.preventDefault();
  const type = $("vehicle-type").value; const source = $("source").value.trim(); const destination = $("destination").value.trim();
  fetch(`${API}/search`, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ vehicleType: type, source, destination }) })
    .then((response) => response.json().then((data) => ({ ok: response.ok, data })))
    .then(({ ok, data }) => {
      if (!ok) throw new Error(data.error || "Search failed");
      const matches = data.vehicles || [];
      window.searchResults = matches;
      $("results-title").textContent = `${source} to ${destination}`;
      $("results-section").classList.remove("hidden"); renderRoutes(matches);
      $("results-section").scrollIntoView({ behavior: "smooth", block: "start" });
    }).catch((error) => toast(error.message));
});
$("booking-form").addEventListener("submit", (event) => {
  event.preventDefault(); if (!$("selected-seat").value) return toast("Select a seat before confirming.");
  fetch(`${API}/book`, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ passengerName: $("passenger-name").value, email: $("passenger-email").value, vehicleId: selectedRoute.id, seatId: Number($("selected-seat").value) }) })
    .then((response) => response.json().then((data) => ({ ok: response.ok, data })))
    .then(({ ok, data }) => { if (!ok || !data.success) throw new Error(data.message || "Booking failed"); $("booking-dialog").close(); $("booking-form").reset(); toast(data.message); $("pnr-input").value = data.pnr; lookupBooking(data.pnr); })
    .catch((error) => toast(error.message));
});
$("registration-form").addEventListener("submit", (event) => {
  event.preventDefault();
  fetch(`${API}/register`, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ name: $("registration-name").value.trim(), email: $("registration-email").value.trim() }) })
    .then((response) => response.json().then((data) => ({ ok: response.ok, data })))
    .then(({ ok, data }) => { if (!ok) throw new Error(data.error || "Registration failed"); $("registration-result").classList.remove("hidden"); $("registration-result").textContent = `Passenger registered in MySQL: ${data.name} (${data.email})`; toast("Passenger saved to the database."); })
    .catch((error) => toast(error.message));
});
function showBooking(data) { const booking = data.booking; $("booking-result").classList.remove("hidden"); $("booking-result").innerHTML = `<strong>${booking.pnr}</strong> · ${booking.vehicle_name} · seat ${booking.seat_number} · ${booking.passenger_name} · <b>${booking.status}</b> <button class="secondary-button cancel-booking">Cancel</button>`; $("booking-result").querySelector(".cancel-booking").addEventListener("click", () => fetch(`${API}/booking?pnr=${encodeURIComponent(booking.pnr)}`, { method: "DELETE" }).then((response) => response.json()).then((result) => { if (!result.success) throw new Error(result.message); showBooking({ booking: { ...booking, status: "CANCELLED" } }); toast("Booking cancelled in the database."); }).catch((error) => toast(error.message))); }
function lookupBooking(pnr) { fetch(`${API}/booking?pnr=${encodeURIComponent(pnr)}`).then((response) => response.json().then((data) => ({ ok: response.ok, data }))).then(({ ok, data }) => { if (!ok || !data.found) throw new Error(data.message || "Booking not found."); showBooking(data); }).catch((error) => toast(error.message)); }
$("lookup-form").addEventListener("submit", (event) => { event.preventDefault(); lookupBooking($("pnr-input").value.trim()); });
$("close-dialog").addEventListener("click", () => $("booking-dialog").close());
