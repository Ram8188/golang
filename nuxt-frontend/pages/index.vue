<template>
  <div>
    <h1>Login</h1>
    <form @submit.prevent="handleLogin">
      <div>
        <label>Email:</label>
        <input type="email" v-model="email" required />
      </div>
      <div>
        <label>Password:</label>
        <input type="password" v-model="password" required />
      </div>
      <button type="submit">Login</button>
    </form>
    <p v-if="error" style="color: red">{{ error }}</p>
  </div>
</template>

<script>
export default {
  data() {
    return {
      email: '',
      password: '',
      error: ''
    };
  },
  methods: {
    async handleLogin() {
        try {
            const response = await fetch("http://localhost:8080/api/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    email: this.email,
                    password: this.password
                })
            });

            const data = await response.json();

            // Check if the login was successful
            if (data.success) {
                // Store the token in localStorage (if needed) or directly store session info
                localStorage.setItem('token', data.token);
                alert("Login successful!");
                this.$router.push("/dashboard"); // Redirect to the dashboard
            } else {
                this.error = "Invalid email or password";
            }
        } catch (error) {
            console.error("Login error:", error); // Log the error for debugging
            this.error = "An error occurred. Please try again.";
        }
    }
}


};
</script>

<style scoped>
label {
  display: block;
  margin: 5px 0;
}
input {
  padding: 8px;
  margin-bottom: 10px;
  width: 100%;
}
button {
  padding: 10px 15px;
}
</style>
