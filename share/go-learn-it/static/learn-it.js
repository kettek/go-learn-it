let LearnIt = {
  editor: {
    textAreaElement: null,
    originalValue: "",
    codeMirror: null,
  },
  output: {
    controlsElement: null,
    outputElement: null,
    btnReset: null,
    btnClear: null,
    btnRun: null,
    result: text => {
      let el = document.createElement('code')
      el.className = "result"
      el.innerText = text
      LearnIt.output.outputElement.appendChild(el)
      LearnIt.output.outputElement.scrollTop = LearnIt.output.outputElement.lastChild.offsetTop
    },
    write: text => {
      let el = document.createElement('code')
      el.innerText = text
      LearnIt.output.outputElement.appendChild(el)
      LearnIt.output.outputElement.scrollTop = LearnIt.output.outputElement.lastChild.offsetTop
    },
    error: text => {
      let el = document.createElement('code')
      el.className = "error"
      el.innerText = text
      LearnIt.output.outputElement.appendChild(el)
      LearnIt.output.outputElement.scrollTop = LearnIt.output.outputElement.lastChild.offsetTop
    },
  },
  canGo: false,
  setup: () => {
    LearnIt.setupEditor()
    LearnIt.setupOutput()
    LearnIt.emit("ready")
  },
  setupEditor: () => {
    LearnIt.editor.textArea = document.getElementById("Editor")
    LearnIt.editor.originalValue = LearnIt.editor.textArea.value
    LearnIt.editor.codeMirror = CodeMirror.fromTextArea(LearnIt.editor.textArea, {
        lineNumbers: true,
        lineWrapping: true,
    })
    LearnIt.editor.codeMirror.on("change", function(cm, change) {
        LearnIt.editor.textArea.value = cm.getValue()
    })
  },
  setupOutput: () => {
    LearnIt.output.controlsElement = document.getElementById("Controls")
    LearnIt.output.outputElement = document.getElementById("OutputPre")
    // Setup our controls
    var btnReset = LearnIt.output.btnReset = document.createElement("a")
    btnReset.className = "button"
    btnReset.id = "ResetButton"
    btnReset.innerText = "Reset"
    btnReset.style.cursor = "not-allowed"
    var btnClear = LearnIt.output.btnClear = document.createElement("a")
    btnClear.className = "button"
    btnClear.id = "ClearButton"
    btnClear.innerText = "Clear"
    btnClear.disabled = true
    btnClear.style.cursor = "not-allowed"
    var btnRun = LearnIt.output.btnRun = document.createElement("a")
    btnRun.className = "button"
    btnRun.id = "RunButton"
    btnRun.innerText = "Run"
    btnRun.style.cursor = "not-allowed"

    btnReset.addEventListener('click', function(e) {
        if (!LearnIt.canGo) return
        LearnIt.editor.codeMirror.setValue(LearnIt.editor.originalValue)
        LearnIt.emit("reset", LearnIt.editor.originalValue)
    })

    btnClear.addEventListener('click', function(e) {
        if (!LearnIt.canGo) return
        LearnIt.output.outputElement.innerHTML = ""
        LearnIt.emit("clear")
    })

    btnRun.addEventListener('click', function(e) {
        if (!LearnIt.canGo) return
        LearnIt.emit("run", LearnIt.editor.codeMirror.getValue())
    })
    LearnIt.output.controlsElement.appendChild(btnReset)
    LearnIt.output.controlsElement.appendChild(btnClear)
    LearnIt.output.controlsElement.appendChild(btnRun)
  }
}
LearnIt.go = function() {
  LearnIt.output.btnReset.style.cursor = null
  LearnIt.output.btnClear.style.cursor = null
  LearnIt.output.btnRun.style.cursor = null
  LearnIt.canGo = true
}
LearnIt._listeners = {}
LearnIt.emit = function(evt, payload) {
  if (!LearnIt._listeners[evt]) { return }
  for (let i in LearnIt._listeners[evt]) {
    LearnIt._listeners[evt][i](payload)
  }
}
LearnIt.on = function(evt, cb) {
  if (!LearnIt._listeners[evt]) {
    LearnIt._listeners[evt] = []
  }
  LearnIt._listeners[evt].push(cb)
}
LearnIt.off = function(evt, cb) {
  if (!LearnIt._listeners[evt]) { return }
  for (let i = 0; i < LearnIt._listeners[evt].length; i++) {
    if (LearnIt._listeners[evt][i] == cb) {
      LearnIt._listeners[evt].splice(i, 1)
    }
  }
}

window.addEventListener("DOMContentLoaded", () => {
  LearnIt.setup()
})