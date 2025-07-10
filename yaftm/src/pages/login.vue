<script lang="ts" setup>
  import FormLogin from '@/components/LoginForm.vue'
</script>

<template>
  <div class="login">
    <h1>ETM Login</h1>
    <p>Welcome to the ETM login page. Please enter your username and password to login.</p>

    <form-login />
    <p>No account yet ? Signup <router-link to="/signup">Here</router-link></p>
  </div>
</template>
