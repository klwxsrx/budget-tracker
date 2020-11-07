<template>
  <b-field
    grouped
    position="is-centered"
  >
    <div class="control">
      <b-button
        :size="controlsSize"
        icon-left="chevron-left"
        icon-pack="fas"
        type="is-dark"
      />
    </div>
    <div class="control">
      <b-button
        :size="controlsSize"
        type="is-dark"
      >
        02 сент. - 02 нояб.
      </b-button>
    </div>
    <div class="control">
      <b-button
        :size="controlsSize"
        icon-left="chevron-right"
        icon-pack="fas"
        type="is-dark"
      />
    </div>
  </b-field>
</template>

<script>
export default { // TODO: real data
  name: 'DateRangePicker',
  computed: {
    controlsSize() {
      return this.$mq === 'tablet' ? 'is-default' : 'is-small'
    },
  },
}
</script>

<style scoped>

</style>