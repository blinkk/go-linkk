const request = require('superagent')

export default class CreateForm {
  constructor(config) {
    this.config = config || {}
    this.form = document.querySelector('#create-form')
    this.status = this.form.querySelector('.create-form__status')
    this.fields = {
      'path': this.form.querySelector('[name=path]'),
      'url': this.form.querySelector('[name=url]'),
      'comment': this.form.querySelector('[name=comment]'),
    }
    this.errors = {}

    this.form.addEventListener('submit', this.processForm.bind(this))
  }

  clearErrors() {
    for (var prop in this.errors) {
      if (this.errors.hasOwnProperty(prop)) {
        delete this.errors[prop]
      }
    }
    return false
  }

  displayErrors() {
    for (var prop in this.fields) {
      if (this.fields.hasOwnProperty(prop)) {
        const field_error = this.fields[prop].parentNode.querySelector('.create-form__errors')
        if (this.errors[prop]) {
          field_error.textContent = this.errors[prop]
        } else {
          field_error.textContent = ''
        }
      }
    }
    return false
  }

  hasErrors() {
    for (var prop in this.errors) {
      if (this.errors.hasOwnProperty(prop)) {
        if (this.errors[prop]) {
          return true
        }
      }
    }
    return false
  }

  processForm(e) {
    if (e.preventDefault) {
      e.preventDefault()
    }

    var values = {
      'path': this.fields['path'].value,
      'url': this.fields['url'].value,
      'comment': this.fields['comment'].value,
    }

    values = this.cleanAndValidate(values)
    this.displayErrors()

    if (this.hasErrors()) {
      console.log('Errors:', this.errors)
      return
    }

    this.status.textContent = 'Submitting...'

    request.post('/_/api/create')
      .type('form')
      .send(values)
      .end(function(err, response) {
        if (err) {
          if (response.body && response.body['errors']) {
            this.status.textContent = 'Please correct the errors and try again.'
            const errors = response.body['errors']
            for (var prop in errors) {
              if (errors.hasOwnProperty(prop)) {
                this.errors[prop.toLowerCase()] = errors[prop]
              }
            }
            this.displayErrors()
          } else {
            console.log(err, response.body)
          }
          return
        }

        console.log(response.body);
        this.status.textContent = 'Success!'
      }.bind(this))
  }

  cleanAndValidate(values) {
    // Cleanup the whitespace.
    for (var prop in values) {
      if (values.hasOwnProperty(prop)) {
        values[prop] = values[prop].trim()
      }
    }

    // Paths cannot end in `/`.
    if (values['path'].endsWith('/')) {
      values['path'] = values['path'].substring(0, values['path'].length - 1)
    }

    // Paths have to start in `/`.
    if (!values['path'].startsWith('/')) {
      values['path'] = '/' + values['path']
    }

    // VALIDATE!
    this.clearErrors()

    return values
  }
}
