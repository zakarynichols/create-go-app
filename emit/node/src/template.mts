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
    <form id="school-form">
      <div class="form-group">
        <label for="schoolId">School ID</label>
        <input type="text" id="schoolId" name="schoolId" required />
      </div>
      <button type="submit">Get School</button>
    </form>
    <div class="school-data">
      <div class="field">
        <label>School ID:</label>
        <span id="school-id"></span>
      </div>
      <div class="field">
        <label>School Name:</label>
        <span id="school-name"></span>
      </div>
      <div class="field">
        <label>Location:</label>
        <span id="school-location"></span>
      </div>
      <div class="field">
        <label>Type:</label>
        <span id="school-type"></span>
      </div>
    </div>
    <script>
      const schoolFormEl = document.querySelector("#school-form");
      const schoolIdInputEl = document.querySelector("#schoolId");
      const schoolIdEl = document.querySelector("#school-id");
      const schoolNameEl = document.querySelector("#school-name");
      const schoolLocation = document.querySelector("#school-location");
      const schoolTypeEl = document.querySelector("#school-type");
      const errorEl = document.querySelector("#error");

      schoolFormEl.addEventListener("submit", async (event) => {
        event.preventDefault();
        const schoolId = schoolIdInputEl.value;
        try {
          const school = await getSchoolById(schoolId);
          schoolIdEl.textContent = school.id;
          schoolNameEl.textContent = school.name;
          schoolLocation.textContent = school.location;
          schoolTypeEl.textContent = school.type;
        } catch (error) {
          console.error(error);
          errorEl.textContent = "failed";
        }
      });

      async function getSchoolById(schoolId) {
        const response = await fetch(\`${serverURL}/schools/\${schoolId}\`, {
          method: "GET",
        });
        if (response.ok) {
          return await response.json();
        }
        if (response.status === 400) {
          throw new Error("Invalid ID supplied");
        }
        if (response.status === 404) {
          throw new Error("School not found");
        }
      }
    </script>
  </body>
</html>

`;
