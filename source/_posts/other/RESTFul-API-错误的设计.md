---
title: RESTFul API 错误的设计
date: 2017-05-13 23:15:53
tags: 
	- API
	- RESTFul
	- 错误
	- 异常
---


不管是WEB还是移动端，如果需要调用Http接口，那么不可避免要处理各种错误，包括客户端参数完整性校验，类型校验，系统错误等。如果一个设计规范的接口错误返回值，不但可以规范调用方统一处理方式，给出更合理的提示，而且后端也有着接口返回规范的作用。这里讨论的其实并不单纯是返回的错误设计，也包括业务上的错误设计。


# 接口返回设计
接口返回需要有http status，错误说明，并且要提供一个完整的错误列表，可以通过简单的错误说明或者错误代码查阅到详细的错误原因。简要错误说明可以只包含简短的文字描述，可以包含错误代码。如果说哪种方式更为合理，我更倾向于错误代码+简要错误说明。如果只有简要错误说明，这个错误描述可能有些描述不合理，如果后端改了说明，那么对应的错误映射也要更改，调用方也要跟着修改。如果只有错误代码，调用方难以一眼看出错误原因，需要查阅错误码，这个比较啰嗦。如果同时提供错误码和简要错误说明，调用方可以通过简要说明快速定位问题，如果要处理错误逻辑，根据错误码表映射即可，这样即使描述改变了，但是错误码是永久不变的，不会影响调用方的处理逻辑。返回的结构如下：
```http
HTTP/1.1 200 OK
Content-Type: application/json;charset=UTF-8

{
	"msg": "错误说明"
}
```

# 服务端错误设计和处理
根据RESTful接口的http status，对于错误，可以分为以下几大类（常用）：
1. 权限错误 401
2. 参数错误 400
3. 重复，冲突 409
4. 找不到资源 404
5. 服务器错误 500

对于以上5中分类，除了第5种系统服务器错误外，其余四种我们可以分别设计一个错误（异常）类型。至于第5中服务器错误为什么不需要，因为我们需要根据上面四种错误是需要单独判断的，如果不在上面的类型中，直接返回500即可。对于上面的错误，在返回或者抛出的时候，应该具体到原因，比如： 参数错误，要写明哪里错误，是手机号不合法还是邮箱不合法，还是哪个参数类型不对，这样便于调用方定位问题。

以Java为例，我们可以定义以下异常：
```java
public class AuthExcepption extends RuntimeException {
	protected int code;
	public AuthExcepption(int code, String msg) {
		this.code = code;
		super(msg);
	}

	public int getCode() {
		return code;
	}
}

// 对于其他三个也是这个格式。
// IllegalArgumentException
// DuplicatedException
// NotFoundException
```

在做异常统一处理的时候，可以判断极少的类型即可，下面是 spring mvc 的统一异常处理方式：
```java
class Body {
	private int code;
	private String msg;
	// setter and getter

	public Body(int code, String msg) {
		this.code = code;
		this.msg = msg;
	}
}

@ExceptionHandler(value = {AuthExcepption.class})
public ResponseEntity<Object> handleAuthExcepption(AuthExcepption ex, WebRequest request) {
	// 记录日志
	return new ResponseEntity(new Body(ex.getCode(),ex.getMessage()),HttpStatus.UNAUTHORIZED);
}

@ExceptionHandler(value = {IllegalArgumentException.class})
public ResponseEntity<Object> handleIllegalArgumentException(IllegalArgumentException ex, WebRequest request) {
	// 记录日志
	return new ResponseEntity(new Body(ex.getCode(),ex.getMessage()), HttpStatus.BAD_REQUEST);
}

@ExceptionHandler(value = {DuplicatedException.class})
public ResponseEntity<Object> handleIllegalArgumentException(DuplicatedException ex, WebRequest request) {
	// 记录日志
	return new ResponseEntity(new Body(ex.getCode(),ex.getMessage()), HttpStatus.CONFLICT);
}

@ExceptionHandler(value = {NotFoundException.class})
public ResponseEntity<Object> handleIllegalArgumentException(NotFoundException ex, WebRequest request) {
	// 记录日志
	return new ResponseEntity(new Body(ex.getCode(),ex.getMessage()),HttpStatus.NOT_FOUND);
}
```


在业务层那么需要这么处理，抛出一种异常，附带具体的错误码和简要消息：
```java
public void addUser(User user) {
	if (!isValid(user.getPhone())) {
		throw new IllegalArgumentException(10001,"手机："+user.getPhone()+" 不合法");
	}
}
```

在上面的业务处理示例中，我们采用了传递错误码的形式，另一种处理方式是定义一个该错误码对应的异常，比如：
```java
public class PhoneIllegalException extends IllegalArgumentException{
	public PhoneIllegalException (String msg) {
		super(10001,msg);
	}
}
```
然后业务层就可以直接抛出此异常了：
```java
public void addUser(User user) {
	if (!isValid(user.getPhone())) {
		throw new PhoneIllegalException("手机："+user.getPhone()+" 不合法");
	}
}
```