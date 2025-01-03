package testutils

import (
	"fmt"
	"net/http"
)

var errorResponseMap = map[ErrorResponseType]struct {
	HTTPStatusCode int
	ResponseBody   string
}{
	ConnectionAborted: {},
	ConnectionRefused: {},
	ConnectionReset:   {},
	ConnectionTimeout: {},

	// Response bodies from https://jira.nyu.edu/browse/DLFA-251, when available.
	// Those values not from https://jira.nyu.edu/browse/DLFA-251 are in square brackets.
	HTTP400BadRequest: {
		HTTPStatusCode: http.StatusBadRequest,
		ResponseBody: makeSolrErrorJSONResponseBody(http.StatusBadRequest,
			"missing content stream", ""),
	},
	HTTP401Unauthorized: {
		HTTPStatusCode: http.StatusUnauthorized,
		ResponseBody: makeSolrErrorJSONResponseBody(http.StatusUnauthorized,
			"["+string(HTTP401Unauthorized)+"]", ""),
	},
	HTTP403Forbidden: {
		HTTPStatusCode: http.StatusForbidden,
		ResponseBody: makeSolrErrorJSONResponseBody(http.StatusForbidden,
			"["+string(HTTP403Forbidden)+"]", ""),
	},
	HTTP404NotFound: {
		HTTPStatusCode: http.StatusNotFound,
		ResponseBody: makeJettyErrorHTMLResponseBody("Error 404 Not Found",
			"HTTP ERROR 404", "Problem accessing /solr/nonexistent-path.",
			"Not Found"),
	},
	HTTP405HTTPMethodNotAllowed: {
		HTTPStatusCode: http.StatusMethodNotAllowed,
		ResponseBody: makeJettyErrorHTMLResponseBody("Error 405 HTTP method POST is not supported by this URL",
			"HTTP ERROR 405", "Problem accessing /solr/admin.html.",
			"HTTP method POST is not supported by this URL"),
	},
	HTTP408RequestTimeout: {
		HTTPStatusCode: http.StatusRequestTimeout,
		ResponseBody: makeSolrErrorJSONResponseBody(http.StatusRequestTimeout,
			"["+string(HTTP408RequestTimeout)+"]", ""),
	},
	HTTP500InternalServerError: {
		HTTPStatusCode: http.StatusInternalServerError,
		ResponseBody: makeSolrErrorJSONResponseBody(http.StatusInternalServerError,
			"", "java.lang.NullPointerException\\n\\tat org.apache.lucene.search.BooleanClause.hashCode(BooleanClause.java:99)\\n\\tat java.util.AbstractList.hashCode(AbstractList.java:541)\\n\\tat org.apache.lucene.search.BooleanQuery.hashCode(BooleanQuery.java:656)\\n\\tat java.util.HashMap.hash(HashMap.java:339)\\n\\tat java.util.HashMap.put(HashMap.java:612)\\n\\tat org.apache.lucene.index.BufferedUpdates.addQuery(BufferedUpdates.java:198)\\n\\tat org.apache.lucene.index.DocumentsWriterDeleteQueue$QueryArrayNode.apply(DocumentsWriterDeleteQueue.java:369)\\n\\tat org.apache.lucene.index.DocumentsWriterDeleteQueue$DeleteSlice.apply(DocumentsWriterDeleteQueue.java:284)\\n\\tat org.apache.lucene.index.DocumentsWriterDeleteQueue.tryApplyGlobalSlice(DocumentsWriterDeleteQueue.java:204)\\n\\tat org.apache.lucene.index.DocumentsWriterDeleteQueue.addDelete(DocumentsWriterDeleteQueue.java:106)\\n\\tat org.apache.lucene.index.DocumentsWriter.deleteQueries(DocumentsWriter.java:141)\\n\\tat org.apache.lucene.index.IndexWriter.deleteDocuments(IndexWriter.java:1465)\\n\\tat org.apache.solr.update.DirectUpdateHandler2.deleteByQuery(DirectUpdateHandler2.java:409)\\n\\tat org.apache.solr.update.processor.RunUpdateProcessor.processDelete(RunUpdateProcessorFactory.java:80)\\n\\tat org.apache.solr.update.processor.UpdateRequestProcessor.processDelete(UpdateRequestProcessor.java:55)\\n\\tat org.apache.solr.update.processor.DistributedUpdateProcessor.doLocalDelete(DistributedUpdateProcessor.java:931)\\n\\tat org.apache.solr.update.processor.DistributedUpdateProcessor.doDeleteByQuery(DistributedUpdateProcessor.java:1433)\\n\\tat org.apache.solr.update.processor.DistributedUpdateProcessor.processDelete(DistributedUpdateProcessor.java:1226)\\n\\tat org.apache.solr.update.processor.LogUpdateProcessor.processDelete(LogUpdateProcessorFactory.java:121)\\n\\tat org.apache.solr.handler.loader.XMLLoader.processDelete(XMLLoader.java:349)\\n\\tat org.apache.solr.handler.loader.XMLLoader.processUpdate(XMLLoader.java:278)\\n\\tat org.apache.solr.handler.loader.XMLLoader.load(XMLLoader.java:174)\\n\\tat org.apache.solr.handler.UpdateRequestHandler$1.load(UpdateRequestHandler.java:99)\\n\\tat org.apache.solr.handler.ContentStreamHandlerBase.handleRequestBody(ContentStreamHandlerBase.java:74)\\n\\tat org.apache.solr.handler.RequestHandlerBase.handleRequest(RequestHandlerBase.java:135)\\n\\tat org.apache.solr.core.SolrCore.execute(SolrCore.java:1967)\\n\\tat org.apache.solr.servlet.SolrDispatchFilter.execute(SolrDispatchFilter.java:777)\\n\\tat org.apache.solr.servlet.SolrDispatchFilter.doFilter(SolrDispatchFilter.java:418)\\n\\tat org.apache.solr.servlet.SolrDispatchFilter.doFilter(SolrDispatchFilter.java:207)\\n\\tat org.eclipse.jetty.servlet.ServletHandler$CachedChain.doFilter(ServletHandler.java:1419)\\n\\tat org.eclipse.jetty.servlet.ServletHandler.doHandle(ServletHandler.java:455)\\n\\tat org.eclipse.jetty.server.handler.ScopedHandler.handle(ScopedHandler.java:137)\\n\\tat org.eclipse.jetty.security.SecurityHandler.handle(SecurityHandler.java:557)\\n\\tat org.eclipse.jetty.server.session.SessionHandler.doHandle(SessionHandler.java:231)\\n\\tat org.eclipse.jetty.server.handler.ContextHandler.doHandle(ContextHandler.java:1075)\\n\\tat org.eclipse.jetty.servlet.ServletHandler.doScope(ServletHandler.java:384)\\n\\tat org.eclipse.jetty.server.session.SessionHandler.doScope(SessionHandler.java:193)\\n\\tat org.eclipse.jetty.server.handler.ContextHandler.doScope(ContextHandler.java:1009)\\n\\tat org.eclipse.jetty.server.handler.ScopedHandler.handle(ScopedHandler.java:135)\\n\\tat org.eclipse.jetty.server.handler.ContextHandlerCollection.handle(ContextHandlerCollection.java:255)\\n\\tat org.eclipse.jetty.server.handler.HandlerCollection.handle(HandlerCollection.java:154)\\n\\tat org.eclipse.jetty.server.handler.HandlerWrapper.handle(HandlerWrapper.java:116)\\n\\tat org.eclipse.jetty.server.Server.handle(Server.java:368)\\n\\tat org.eclipse.jetty.server.AbstractHttpConnection.handleRequest(AbstractHttpConnection.java:489)\\n\\tat org.eclipse.jetty.server.BlockingHttpConnection.handleRequest(BlockingHttpConnection.java:53)\\n\\tat org.eclipse.jetty.server.AbstractHttpConnection.content(AbstractHttpConnection.java:953)\\n\\tat org.eclipse.jetty.server.AbstractHttpConnection$RequestHandler.content(AbstractHttpConnection.java:1014)\\n\\tat org.eclipse.jetty.http.HttpParser.parseNext(HttpParser.java:861)\\n\\tat org.eclipse.jetty.http.HttpParser.parseAvailable(HttpParser.java:240)\\n\\tat org.eclipse.jetty.server.BlockingHttpConnection.handle(BlockingHttpConnection.java:72)\\n\\tat org.eclipse.jetty.server.bio.SocketConnector$ConnectorEndPoint.run(SocketConnector.java:264)\\n\\tat org.eclipse.jetty.util.thread.QueuedThreadPool.runJob(QueuedThreadPool.java:608)\\n\\tat org.eclipse.jetty.util.thread.QueuedThreadPool$3.run(QueuedThreadPool.java:543)\\n\\tat java.lang.Thread.run(Thread.java:748)\\n"),
	},
	HTTP502BadGateway: {
		HTTPStatusCode: http.StatusBadGateway,
		ResponseBody: makeJettyErrorHTMLResponseBody("Error 502 Bad Gateway",
			"HTTP ERROR 502", "["+string(HTTP502BadGateway)+"]",
			"Bad Gateway"),
	},
	HTTP503ServiceUnavailable: {
		HTTPStatusCode: http.StatusServiceUnavailable,
		ResponseBody: makeJettyErrorHTMLResponseBody("Error 503 Service Unavailable",
			"HTTP ERROR 503", "["+string(HTTP503ServiceUnavailable)+"]",
			"Service Unavailable"),
	},
	HTTP504GatewayTimeout: {
		HTTPStatusCode: http.StatusGatewayTimeout,
		ResponseBody: makeJettyErrorHTMLResponseBody("Error 504 Gateway Timeout",
			"HTTP ERROR 504", "["+string(HTTP504GatewayTimeout)+"]",
			"Gateway Timeout"),
	},

	ConnectionTimeoutPermanent: {},
}

func makeSolrErrorJSONResponseBody(code int, msg string, trace string) string {
	var key, value string

	if msg != "" {
		key = "msg"
		value = msg
	} else if trace != "" {
		key = "trace"
		value = trace
	} else {
		// Should never get here
		panic("makeSolrErrorJSONResponseBody() helper called with no `msg` or `trace`")
	}

	return fmt.Sprintf(`{"responseHeader":{"status":%d,"QTime":0},"error":{"%s":"%s","code":%d}}`,
		code, key, value, code)
}

func makeJettyErrorHTMLResponseBody(title string, heading string, problem string,
	reason string) string {
	return fmt.Sprintf(
		`<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=ISO-8859-1"/>
<title>%s</title>
</head>
<body><h2>%s</h2>
<p>%s. Reason:
<pre>    %s</pre></p><hr /><i><small>Powered by Jetty://</small></i><br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                
<br/>                                                

</body>
</html>`, title, heading, problem, reason)
}
