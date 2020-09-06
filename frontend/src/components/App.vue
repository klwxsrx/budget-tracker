<template>
  <div class="some_style">
    {{ message }}
  </div>
</template>

<script>
export default {
  data() {
    return {
      message: "Hello, Vue!"
    };
  }
}
</script>

<style lang="scss" scoped>
  .some_style
  {
    background-color: purple;
  }
</style>