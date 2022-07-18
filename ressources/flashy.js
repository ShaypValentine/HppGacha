var Flashy = (function () {
  'use strict';

  function _classCallCheck(instance, Constructor) {
    if (!(instance instanceof Constructor)) {
      throw new TypeError("Cannot call a class as a function");
    }
  }

  function _defineProperties(target, props) {
    for (var i = 0; i < props.length; i++) {
      var descriptor = props[i];
      descriptor.enumerable = descriptor.enumerable || false;
      descriptor.configurable = true;
      if ("value" in descriptor) descriptor.writable = true;
      Object.defineProperty(target, descriptor.key, descriptor);
    }
  }

  function _createClass(Constructor, protoProps, staticProps) {
    if (protoProps) _defineProperties(Constructor.prototype, protoProps);
    if (staticProps) _defineProperties(Constructor, staticProps);
    return Constructor;
  }

  function _inherits(subClass, superClass) {
    if (typeof superClass !== "function" && superClass !== null) {
      throw new TypeError("Super expression must either be null or a function");
    }

    subClass.prototype = Object.create(superClass && superClass.prototype, {
      constructor: {
        value: subClass,
        writable: true,
        configurable: true
      }
    });
    if (superClass) _setPrototypeOf(subClass, superClass);
  }

  function _getPrototypeOf(o) {
    _getPrototypeOf = Object.setPrototypeOf ? Object.getPrototypeOf : function _getPrototypeOf(o) {
      return o.__proto__ || Object.getPrototypeOf(o);
    };
    return _getPrototypeOf(o);
  }

  function _setPrototypeOf(o, p) {
    _setPrototypeOf = Object.setPrototypeOf || function _setPrototypeOf(o, p) {
      o.__proto__ = p;
      return o;
    };

    return _setPrototypeOf(o, p);
  }

  function isNativeReflectConstruct() {
    if (typeof Reflect === "undefined" || !Reflect.construct) return false;
    if (Reflect.construct.sham) return false;
    if (typeof Proxy === "function") return true;

    try {
      Date.prototype.toString.call(Reflect.construct(Date, [], function () {}));
      return true;
    } catch (e) {
      return false;
    }
  }

  function _construct(Parent, args, Class) {
    if (isNativeReflectConstruct()) {
      _construct = Reflect.construct;
    } else {
      _construct = function _construct(Parent, args, Class) {
        var a = [null];
        a.push.apply(a, args);
        var Constructor = Function.bind.apply(Parent, a);
        var instance = new Constructor();
        if (Class) _setPrototypeOf(instance, Class.prototype);
        return instance;
      };
    }

    return _construct.apply(null, arguments);
  }

  function _isNativeFunction(fn) {
    return Function.toString.call(fn).indexOf("[native code]") !== -1;
  }

  function _wrapNativeSuper(Class) {
    var _cache = typeof Map === "function" ? new Map() : undefined;

    _wrapNativeSuper = function _wrapNativeSuper(Class) {
      if (Class === null || !_isNativeFunction(Class)) return Class;

      if (typeof Class !== "function") {
        throw new TypeError("Super expression must either be null or a function");
      }

      if (typeof _cache !== "undefined") {
        if (_cache.has(Class)) return _cache.get(Class);

        _cache.set(Class, Wrapper);
      }

      function Wrapper() {
        return _construct(Class, arguments, _getPrototypeOf(this).constructor);
      }

      Wrapper.prototype = Object.create(Class.prototype, {
        constructor: {
          value: Wrapper,
          enumerable: false,
          writable: true,
          configurable: true
        }
      });
      return _setPrototypeOf(Wrapper, Class);
    };

    return _wrapNativeSuper(Class);
  }

  function _assertThisInitialized(self) {
    if (self === void 0) {
      throw new ReferenceError("this hasn't been initialised - super() hasn't been called");
    }

    return self;
  }

  function _possibleConstructorReturn(self, call) {
    if (call && (typeof call === "object" || typeof call === "function")) {
      return call;
    }

    return _assertThisInitialized(self);
  }

  function _slicedToArray(arr, i) {
    return _arrayWithHoles(arr) || _iterableToArrayLimit(arr, i) || _nonIterableRest();
  }

  function _arrayWithHoles(arr) {
    if (Array.isArray(arr)) return arr;
  }

  function _iterableToArrayLimit(arr, i) {
    var _arr = [];
    var _n = true;
    var _d = false;
    var _e = undefined;

    try {
      for (var _i = arr[Symbol.iterator](), _s; !(_n = (_s = _i.next()).done); _n = true) {
        _arr.push(_s.value);

        if (i && _arr.length === i) break;
      }
    } catch (err) {
      _d = true;
      _e = err;
    } finally {
      try {
        if (!_n && _i["return"] != null) _i["return"]();
      } finally {
        if (_d) throw _e;
      }
    }

    return _arr;
  }

  function _nonIterableRest() {
    throw new TypeError("Invalid attempt to destructure non-iterable instance");
  }

  window.onload = function () {
    console.log('flashy');
    var style = document.createElement('style');
    style.textContent = "\n        flash-messages {\n            position: fixed;\n            display: flex;\n            flex-wrap: wrap;\n            left: 10px;\n            bottom: 0px;\n            width: 17%;\n            font-family: sans-serif;\n            min-width: 300px;\n            min-height: 100px;\n        }\n    ";
    document.querySelector('body').appendChild(style);
  };

  function triggerAnimation(duration, className, element) {
    element.classList.add(className);
    element.addEventListener('animationend', function () {
      element.classList.remove(className);
    });
  }

  var FlashMessages =
  /*#__PURE__*/
  function (_HTMLElement) {
    _inherits(FlashMessages, _HTMLElement);

    function FlashMessages() {
      var _this;

      _classCallCheck(this, FlashMessages);

      _this = _possibleConstructorReturn(this, _getPrototypeOf(FlashMessages).call(this));
      setTimeout(function () {
        _this.maxMessages = _this.dataset['maxMessages'] || 10;
      }, 0);
      var style = document.createElement('style');
      style.textContent = "\n        .flash-message {\n            width: 100%;\n            height: max-content;\n            margin-bottom: 5%;\n            box-shadow: 0 1px 3px rgba(0, 0, 0, 0.12), 0 1px 2px rgba(0, 0, 0, 0.24);\n            transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);\n            border-radius: 10px;\n            background-color: #ecf0f1;\n            font-family: 'Roboto', sans-serif;\n            display: flex;\n            align-items: center;\n            position: relative;\n            -webkit-box-sizing: border-box; /* Safari/Chrome, other WebKit */\n            -moz-box-sizing: border-box;    /* Firefox, other Gecko */\n            box-sizing: border-box;         /* Opera/IE 8+ */\n            padding-right: 10px;\n            \n        }\n       \n        .flash-message:hover {\n            box-shadow: 0 14px 28px rgba(0, 0, 0, 0.25), 0 10px 10px rgba(0, 0, 0, 0.22);\n        }\n       \n        .flash-message.error {\n            background-color: #ff6060;\n        }\n \n        .flash-message.error > .left-part {\n            background-color: #ff3535;\n        }\n       \n        .flash-message.info {\n            background-color: #59c9f9;\n        }\n \n        .flash-message.info > .left-part {\n            background-color: #32c0ff;\n        }\n       \n        .flash-message.success {\n            background-color: #39c16c;\n        }\n \n        .flash-message.success > .left-part {\n            background-color: #33ad61;\n        }\n       \n        .flash-message.warning {\n            background-color: #ffc744;\n        }\n \n        .flash-message.warning > .left-part {\n            background-color: #ffba1c;\n        }\n       \n        .left-part {\n            height: 100%;\n            width: 25%;\n            display: flex;\n            flex-shrink: 1;\n            align-items: center;\n            justify-content: center;\n            color: white;\n            font-size: 3em;\n            border-radius: 10px 0 0 10px;\n        }\n \n        .right-part {\n            display: inline-block;\n            flex-shrink: 100;\n            padding-left: 10px;\n            padding-top: 10px;\n            padding-bottom: 10px;\n        }\n       \n        .flash-title {\n            display: block;\n            font-size: 1.2em;\n            font-weight: 500;\n            font-family: 'Roboto', sans-serif;\n            color: white;\n        }\n       \n        .flash-msg {\n            margin-top: 3px;\n            display: block;\n            color: white;\n            font-family: 'Roboto', sans-serif;\n            font-weight: 300;\n            font-size: .8em;\n        }\n       \n        .close-flash {\n            box-shadow: 0 1px 3px rgba(0, 0, 0, 0.12), 0 1px 2px rgba(0, 0, 0, 0.24);\n            position: absolute;\n            right: -8px;\n            top: -8px;\n            width: 22px;\n            height: 22px;\n            background-color: #f1f1f1;\n            border-radius: 11px;\n            display: flex;\n            align-items: center;\n            justify-content: center;\n            color: #7f8c8d;\n            font-size: .8em;\n        }\n       \n        .close-flash:hover {\n            cursor: pointer;\n            color: #e74c3c;\n        }\n\n        .flashy-button {\n            color: #1663fc;\n            font-size: .8em;\n            font-align: center;\n            padding-left: 5px;\n            padding-right: 5px;\n        }\n\n        .flashy-button:nth-child(0) {\n            padding-left: 0;\n        }\n\n        .flashy-button:hover {\n            cursor: pointer;\n            text-decoration: underline;\n        }\n\n        .flashy-seperator {\n\n        }\n       \n        .movein {\n            animation-name: movein;\n            animation-duration: .2s;\n        }\n       \n        .movedown {\n            animation-name: movedown;\n            animation-duration: .2s;\n        }\n       \n        .moveup {\n            animation-name: moveup;\n            animation-duration: .2s;\n        }\n       \n        .fade {\n            animation-name: fade;\n            animation-duration: .2s;\n            opacity: 0;\n        }\n       \n        @keyframes moveup {\n            from {\n                transform: translateY(100%);\n            }\n            to {\n                transform: translateY(0%);\n            }\n        }\n       \n        @keyframes movedown {\n            from {\n                transform: translateY(-100%);\n            }\n            to {\n                transform: translateY(0);\n            }\n        }\n       \n        @keyframes movein {\n            from {\n                transform: translateX(-100%);\n            }\n            to {\n                transform: translateX(0);\n            }\n        }\n       \n        @keyframes fade {\n            from {\n                opacity: 1;\n            }\n            to {\n                opacity: 0;\n            }\n        }\n        ";
      _this.shadowRootObj = _this.attachShadow({
        mode: 'open'
      });

      _this.shadowRootObj.appendChild(style);

      return _this;
    }

    _createClass(FlashMessages, [{
      key: "remove_child",
      value: function remove_child(elem) {
        var _this2 = this;

        elem.classList.add('fade');
        setTimeout(function () {
          if (elem.nextElementSibling) {
            triggerAnimation(200, 'moveup', elem.nextElementSibling);
          }

          if (elem.nextElementSibling !== null && elem.nextElementSibling.nextElementSibling) {
            triggerAnimation(200, 'moveup', elem.nextElementSibling.nextElementSibling);
          }

          elem.remove();

          _this2.set_bottom();
        }, 200);
      }
    }, {
      key: "create_child",
      value: function create_child(options) {
        var _this3 = this;

        var emoji = {
          error: '😨',
          warning: '😶',
          success: '😇',
          info: '😜'
        };
        var icon = emoji[options.type];

        if (options.styles) {
          if (options.styles.icon) {
            if (options.styles.icon.type === 'unicode') {
              icon = options.styles.icon.val;
            }
          }
        }

        var newFlash = document.createElement('div');
        newFlash.classList.add('flash-message');
        newFlash.classList.add(options.type);

        if (options.styles && options.styles.flashColor) {
          newFlash.style.backgroundColor = options.styles.flashColor;
        }

        var leftPart = document.createElement('div');
        leftPart.classList.add('left-part');
        leftPart.innerHTML = icon;

        if (options.styles && options.styles.iconBackgroundColor) {
          leftPart.style.backgroundColor = options.styles.iconBackgroundColor;
        }

        newFlash.appendChild(leftPart);
        var rightPart = document.createElement('div');
        rightPart.classList.add('right-part');
        newFlash.appendChild(rightPart);
        var ttl = document.createElement('span');
        ttl.innerHTML = options.title;
        ttl.classList.add('flash-title');

        if (options.styles) {
          if (options.styles.titleTextColor) {
            ttl.style.color = options.styles.titleTextColor;
          }

          if (options.styles.titleTextFont) {
            ttl.style.fontFamily = options.styles.titleTextFont;
          }
        }

        rightPart.appendChild(ttl);
        var msg = document.createElement('span');
        msg.innerHTML = options.message;
        msg.classList.add('flash-msg');

        if (options.styles) {
          if (options.styles.msgTextColor) {
            msg.style.color = options.styles.msgTextColor;
          }

          if (options.styles.msgTextFont) {
            msg.style.fontFamily = options.styles.msgTextFont;
          }
        }

        rightPart.appendChild(msg);

        if (options.buttons) {
          var _iteratorNormalCompletion = true;
          var _didIteratorError = false;
          var _iteratorError = undefined;

          try {
            var _loop = function _loop() {
              var _step$value = _slicedToArray(_step.value, 2),
                  i = _step$value[0],
                  button = _step$value[1];

              var btn = document.createElement('span');
              btn.classList.add('flashy-button');
              btn.innerHTML = button.text;
              btn.addEventListener('click', function (event) {
                if (button.action) {
                  button.action();
                }

                if (button.closesFlash) {
                  _this3.remove_child(btn.parentElement.parentElement);
                }
              });

              if (options.styles) {
                if (options.styles.linkTextColor) {
                  btn.style.color = options.styles.linkTextColor;
                }

                if (options.styles.linkTextFont) {
                  btn.style.fontFamily = options.styles.linkTextFont;
                }
              }

              rightPart.appendChild(btn);

              if (i !== options.buttons.length - 1) {
                var sep = document.createElement('span');
                sep.classList.add('flashy-seperator');
                sep.innerHTML = '·';
                rightPart.appendChild(sep);
              }
            };

            for (var _iterator = options.buttons.entries()[Symbol.iterator](), _step; !(_iteratorNormalCompletion = (_step = _iterator.next()).done); _iteratorNormalCompletion = true) {
              _loop();
            }
          } catch (err) {
            _didIteratorError = true;
            _iteratorError = err;
          } finally {
            try {
              if (!_iteratorNormalCompletion && _iterator["return"] != null) {
                _iterator["return"]();
              }
            } finally {
              if (_didIteratorError) {
                throw _iteratorError;
              }
            }
          }
        }

        if (options.globalClose) {
          var span = document.createElement('div');
          span.innerHTML = 'X';
          span.classList.add('close-flash');
          span.addEventListener('click', function (event) {
            var flashElem = event.srcElement.parentElement || event.target.parentElement;

            _this3.remove_child(flashElem);
          });
          newFlash.appendChild(span);
        }

        return newFlash;
      }
    }, {
      key: "set_bottom",
      value: function set_bottom() {
        this.style.bottom = "".concat(this.shadowRootObj.children[0].clientHeight - this.clientHeight + 20, "px");
        this.shadowRootObj.children[0].style.height = "".concat(this.shadowRootObj.children[0].clientHeight, "px");
      }
    }, {
      key: "add_child",
      value: function add_child(options) {
        var _this4 = this;

        var newFlash = this.create_child(options);
        var length = this.shadowRootObj.children.length;

        if (length - 1 >= this.maxMessages) {
          this.shadowRootObj.lastChild.previousSibling.remove();
        }

        this.shadowRootObj.insertBefore(newFlash, this.shadowRootObj.firstChild);
        this.set_bottom(newFlash);
        triggerAnimation(200, 'movein', newFlash);

        if (newFlash.nextElementSibling) {
          triggerAnimation(200, 'movedown', newFlash.nextElementSibling);
        }

        if (options.expiry > 0) {
          setTimeout(function () {
            return _this4.remove_child(newFlash);
          }, options.expiry);
        }
      }
    }]);

    return FlashMessages;
  }(_wrapNativeSuper(HTMLElement));

  window.customElements.define('flash-messages', FlashMessages);

  Flashy = function Flashy(selector, options) {
    document.querySelector(selector).add_child(options);
  };

  return Flashy;

}());
