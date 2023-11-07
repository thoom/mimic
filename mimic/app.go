package mimic

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type MimicContext string

const (
	ContextJose MimicContext = "mimicJose"
)

func TermsHandler(w http.ResponseWriter, r *http.Request) {
	txt := `TERMS OF SERVICE 

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.`

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(txt))
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add some route logging
		log.Printf("%s %s %s\n", r.Method, r.RequestURI, r.Header.Get("Content-Type"))
		next.ServeHTTP(w, r)
	})
}

func JoseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)

		var joseJson JoseJson
		if r.Header.Get("Content-Type") == "application/jose+json" {
			json.Unmarshal([]byte(buf.Bytes()), &joseJson)
			joseJson.DecodeProtected()
			joseJson.DecodePayload()

			//TODO: validate the JWT
		}

		ctx := context.WithValue(r.Context(), ContextJose, joseJson)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
