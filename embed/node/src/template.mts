const server = {
  host: process.env.GO_HOST,
  port: process.env.GO_PORT,
};

// serverURL is the URL of the Go HTTP server.
const serverURL = `http://${server.host}:${server.port}`; // On the host machine is http://localhost:1111

export const indexHTML = `
<!DOCTYPE html>
<html>
  <head>
    <title>HTTP Requests Example</title>
  </head>
  <body>
    <div id="error"></div>
    <form id="thing-form">
      <div class="form-group">
        <label for="thingId">thing ID</label>
        <input type="text" id="thingId" name="thingId" required />
      </div>
      <button type="submit">Get thing</button>
    </form>
    <div class="thing-data">
      <div class="field">
        <label>thing ID:</label>
        <span id="thing-id"></span>
      </div>
      <div class="field">
        <label>thing Name:</label>
        <span id="thing-name"></span>
      </div>
      <div class="field">
        <label>Location:</label>
        <span id="thing-location"></span>
      </div>
      <div class="field">
        <label>Type:</label>
        <span id="thing-type"></span>
      </div>
    </div>
    <script>
      const thingFormEl = document.querySelector("#thing-form");
      const thingIdInputEl = document.querySelector("#thingId");
      const thingIdEl = document.querySelector("#thing-id");
      const thingNameEl = document.querySelector("#thing-name");
      const thingLocation = document.querySelector("#thing-location");
      const thingTypeEl = document.querySelector("#thing-type");
      const errorEl = document.querySelector("#error");

      thingFormEl.addEventListener("submit", async (event) => {
        event.preventDefault();
        const thingId = thingIdInputEl.value;
        try {
          const thing = await getThingById(thingId);
          thingIdEl.textContent = thing.id;
          thingNameEl.textContent = thing.name;
          thingLocation.textContent = thing.location;
          thingTypeEl.textContent = thing.type;
        } catch (error) {
          console.error(error);
          errorEl.textContent = "failed";
        }
      });

      async function getThingById(thingId) {
        const response = await fetch(\`${serverURL}/things/\${thingId}\`, {
          method: "GET",
        });
        if (response.ok) {
          return await response.json();
        }
        if (response.status === 400) {
          throw new Error("Invalid ID supplied");
        }
        if (response.status === 404) {
          throw new Error("thing not found");
        }
      }
    </script>
  </body>
</html>

`;
