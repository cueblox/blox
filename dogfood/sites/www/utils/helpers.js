export function formatContentfulDate(contentfulDate) {
  const convertedDate = new Date(contentfulDate)

  const year = convertedDate.getFullYear()
  const date = convertedDate.getDate()

  const months = [
    'January',
    'February',
    'March',
    'April',
    'May',
    'June',
    'July',
    'August',
    'September',
    'October',
    'November',
    'December'
  ]

  const monthIndex = convertedDate.getMonth()
  const monthName = months[monthIndex]

  const formatted = `${monthName} ${date}, ${year}`
  return formatted;
}
