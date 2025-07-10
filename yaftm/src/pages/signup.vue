<script setup lang="ts">
  import SignupForm from '@/components/SignUpForm.vue'
</script>

<template>
  <div class="signup">
    <h1>ETM Signup</h1>
    <p>Welcome to the ETM Signup page. Please enter the requested information to create an account.</p>

    <SignupForm />
    <p>Already have an account ? Login <router-link to="/login">Here</router-link></p>
  </div>
</template>

<style scoped lang="sass">

</style>
