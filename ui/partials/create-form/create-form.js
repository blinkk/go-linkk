const request = require('superagent')

export default class CreateForm {
  constructor(config) {
    this.config = config || {}
    this.form = document.querySelector('#create-form')
    this.fields = {
      'path': this.form.querySelector('[name=path]'),
      'url': this.form.querySelector('[name=url]'),
      'comment': this.form.querySelector('[name=comment]'),
    }
    this.errors = {}

    this.form.addEventListener('submit', this.processForm.bind(this))
  }

  hasErrors() {
    for (var prop in self.errors) {
      if (self.errors.hasOwnProperty(prop)) {
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

    if (this.hasErrors()) {
      console.log('Errors:', this.errors)
      return
    }

    request.post('/_/api/create')
      .type('form')
      .send(values)
      .end(function(err, response) {
        if (err) {
          console.error(err);
          return;
        }

        console.log('Success!');
      });
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

    return values
  }
}
