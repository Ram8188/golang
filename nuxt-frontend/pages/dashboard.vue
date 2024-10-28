<template>
    <div>
      <div class="sidebar">
        <div class="logo">
          <img src="" alt="Logo" />
        </div>
        <nav>
          <ul>
            <li><router-link to="/dashboard">Dashboard</router-link></li>
            <li><router-link to="/profile">Profile</router-link></li>
            <li><router-link to="/settings">Settings</router-link></li>
            <li @click="logout">Logout</li>
          </ul>
        </nav>
      </div>
      <div class="main-content">
        <h1>Welcome to the Dashboard!</h1>
        <!-- Other dashboard content goes here -->
      </div>
    </div>
  </template>
  
  <script>
  export default {
    name: "Dashboard",
    created() {
      this.checkSession();
    },
    methods: {
      async checkSession() {
        try {
          const response = await fetch("http://localhost:8080/api/verify-session", {
            method: "GET",
            credentials: "include", // Include cookies in the request
          });
          if (!response.ok) {
            this.$router.push("/"); // Redirect to index page if not logged in
          }
        } catch (error) {
          console.error("Session verification error:", error);
          this.$router.push("/"); // Redirect on error as well
        }
      },
      async logout() {
        // Implement logout functionality
        localStorage.removeItem('token'); // Clear token if applicable
        this.$router.push("/"); // Redirect to index page after logout
      },
    },
  };
  </script>
  
  <style scoped>
  .sidebar {
    width: 200px;
    background-color: #f4f4f4;
    height: 100vh;
    position: fixed;
  }
  
  .logo {
    padding: 20px;
    text-align: center;
  }
  
  .logo img {
    max-width: 100%;
    height: auto;
  }
  
  nav {
    margin-top: 20px;
  }
  
  nav ul {
    list-style: none;
    padding: 0;
  }
  
  nav ul li {
    padding: 10px;
  }
  
  nav ul li:hover {
    background-color: #ddd;
  }
  
  .main-content {
    margin-left: 220px;
    padding: 20px;
  }
  </style>
  