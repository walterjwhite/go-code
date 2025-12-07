document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("contactForm");
  const submitBtn = document.getElementById("submitBtn");
  const statusEl = document.getElementById("status");

  function setStatus(text, cls) {
    statusEl.textContent = text;
    statusEl.className = cls || "";
  }

  form.addEventListener("submit", async (ev) => {
    ev.preventDefault();
    setStatus("", "");
    submitBtn.disabled = true;
    submitBtn.textContent = "Sending...";

    const data = {
      name: form.name.value.trim(),
      email: form.email.value.trim(),
      subject: form.subject.value.trim(),
      message: form.message.value.trim(),
    };

    try {
      const res = await fetch("/api/contact", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      if (res.ok) {
        setStatus("Message sent — thank you!", "success");
        form.reset();
      } else {
        const body = await res.json().catch(() => ({}));
        const err = body && body.error ? body.error : `${res.status} ${res.statusText}`;
        setStatus("Error sending message: " + err, "error");
      }
    } catch (err) {
      setStatus("Network error: " + err.message, "error");
    } finally {
      submitBtn.disabled = false;
      submitBtn.textContent = "Send message";
    }
  });
});