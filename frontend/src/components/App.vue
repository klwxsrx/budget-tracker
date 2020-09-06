<template>
  <div>
    <b-button size="is-small" icon-pack="fab" icon-left="github">
      {{message}}
    </b-button>
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

<style scoped>

</style>