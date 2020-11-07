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
        icon-left="calendar-alt"
        icon-pack="far"
        type="is-dark"
      >
        {{ displayingDate }}
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
import moment from 'moment'

export default {
  name: 'TransactionDateRangeSelector',
  props: {
    minDate: {
      type: Date,
      default: null,
    },
    maxDate: {
      type: Date,
      default: null,
    },
    startDate: {
      type: Date,
      default: null,
    },
    endDate: {
      type: Date,
      default: null,
    },
  },
  computed: {
    displayingDate() {
      if (this.areMonthPeriodSelected(this.$props.startDate, this.$props.endDate)) {
        const date = moment(this.$props.startDate).format('MMMM YYYY')
        return date[0].toUpperCase() + date.slice(1)
      }
      const areTheSameDaySelected = (this.$props.startDate - this.$props.endDate) === 0
      if (areTheSameDaySelected) {
        return moment(this.$props.startDate).format('DD.MM.YYYY')
      }
      return moment(this.$props.startDate).format('DD.MM.YYYY') + ' - ' + moment(this.$props.endDate).format('DD.MM.YYYY')
    },
    controlsSize() {
      return this.$mq === 'tablet' ? 'is-default' : 'is-small'
    },
  },
  methods: {
    // Month period is selected when:
    // 1) startDate and endDate are first and last days of the same month
    // 2) the difference between startDate and endDate is month, e.g. June 10 - July 9, January 31 - February 28
    areMonthPeriodSelected(startDate, endDate) {
      const areTheSameYear = startDate.getFullYear() === endDate.getFullYear()
      const lastDayOfMonthStartDate = new Date(startDate.getFullYear(), startDate.getMonth() + 1, 0).getDate()
      const startDateIsFirstDayOfMonth = startDate.getDate() === 1
      const endDateIsLastDayOfTheSameMonth = startDate.getMonth() === endDate.getMonth() && endDate.getDate() === lastDayOfMonthStartDate

      if (startDateIsFirstDayOfMonth && endDateIsLastDayOfTheSameMonth && areTheSameYear) {
        return true
      }

      const endDateIsInTheNextYear = endDate.getFullYear() - startDate.getFullYear() === 1
      const endDateIsInTheNextMonth = (endDate.getMonth() - startDate.getMonth() === 1 && areTheSameYear) || (endDate.getMonth() - startDate.getMonth() === -11 && endDateIsInTheNextYear)
      const startDateDayIsTheNextDay = startDate.getDate() - endDate.getDate() === 1

      const lastDayOfMonthEndDate = new Date(endDate.getFullYear(), endDate.getMonth() + 1, 0).getDate()
      const endDateIsTheLastDayOfMonth = endDate.getDate() === lastDayOfMonthEndDate
      const startDateDayIsGreaterThanTheEndDate = startDate.getDate() > lastDayOfMonthEndDate

      return endDateIsInTheNextMonth && (startDateDayIsTheNextDay || startDateDayIsGreaterThanTheEndDate && endDateIsTheLastDayOfMonth)
    },
  },
}
</script>

<style scoped>

</style>