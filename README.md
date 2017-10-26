# CMDBc
The _**Configuration Management Database Client**_ or **CMDBc** is a utility that manages information about devices attached to end-user workstations and reports that information to a centralized repository over a RESTful JSON API provided by the complementary server component, the _**Configuration Management Database Daemon**_ or **CMDBd**. Detailed documentation for **CMDBd** is located [here](https://github.com/jscherff/cmdbd/blob/master/README.md). 

**CMDBc** can register or _"check-in"_ attached devices with the server, obtain unique serial numbers from the server for devices that support serial number configuration, perform audits against previous device configurations, and report configuration changes found during the audit to the server for logging and analysis.

### System Requirements
**CMDBc** is written in **Go** and can be compiled for any operating system and architecture. It is intended to be installed on end-user workstations running **Microsoft Windows 7** or higher and should be invoked by a centralized management solution like **IBM BigFix**.

### Installation
To implement **CMDBc** quickly with minimal configuration, please see the [**QuickStart** document.](https://github.com/jscherff/cmdbc/blob/master/QUICKSTART.md)

Pre-compiled binaries are available for both 32- and 64-bit Windows systems and can be installed in any folder along with the required JSON configuration file:

* [**`cmdbc.exe`**](https://github.com/jscherff/cmdbc/raw/master/i686/cmdbc.exe) (32-bit Windows 7 or higher)
* [**`cmdbc.exe`**](https://github.com/jscherff/cmdbc/raw/master/x86_64/cmdbc.exe) (64-bit Windows 7 or higher)
* [**`config.json`**](https://github.com/jscherff/cmdbc/raw/master/config.json) (Configuration file)

### Configuration
The JSON configuration file, [`config.json`](https://github.com/jscherff/cmdbd/blob/master/config.json), is mostly self-explanatory. The default settings are sane and you should not have to change them in most use cases.

#### API Settings
The **API** section of the configuration file contains parameters for communicating with the **CMDBd** server and URL paths for the REST API endpoints.
```json
"API": {
    "Server": "http://cmdbsvcs-prd-01.24hourfit.com:8080",
    "Endpoint": {
        "usbCiCheckinV1": "v1/usbci/checkin",
        "usbCiCheckoutV1": "v1/usbci/checkout",
        "usbCiNewSnV1": "v1/usbci/newsn",
        "usbCiAuditV1": "v1/usbci/audit",
        "usbMetaVendorV1": "v1/usbmeta/vendor",
        "usbMetaProductV1": "v1/usbmeta/product",
        "usbMetaClassV1": "v1/usbmeta/class",
        "usbMetaSubClassV1": "v1/usbmeta/subclass",
        "usbMetaProtocolV1": "v1/usbmeta/protocol"
    }
}
```
* **`Server`** is the base URL for the server hosting the REST API.
* **`Endpoints`** is a collection of URL paths that represent the base of the REST API endpoints on the server. The API endpoints and their parameters are described more fully in the [API Endpoints](https://github.com/jscherff/cmdbd/blob/master/README.md#api-endpoints) section of the server documentation. You should not modify anything in this section unless asked to do so by a systems administrator or application designer.
    * **`v1/usbci/checkin`** is the base path of the API on which the client submits configuration information for a new device or update information for an existing device.
    * **`v1/usbci/checkout`** is the base path of the API on which the client obtains configuration information for a previously-registered, serialized device in order to perform a change audit.
    * **`v1/usbci/newsn`** is the base path of the API on which the client obtains a new unique serial number from the server for assignment to the attached device.
    * **`v1/usbci/audit`** is the base path of the API on which the client submit the results of a change audit on a serialized device. Results include the attribute name, previous value, and new value for each modified attribute.
    * **`v1/usbmeta/vendor`** is the base path of the API on which the client obtains the USB vendor name by providing the vendor ID.
    * **`v1/usbmeta/product`** is the base path of the API on which the client obtains the USB vendor and product names by providing the vendor and product IDs.
    * **`v1/usbmeta/class`** is the base path of the API on which the client obtains the USB class description by providing the class ID.
    * **`v1/usbmeta/subclass`** is the base path of the API on which the client obtains the USB class and subclass descriptions by providing the class and subclass IDs.
    * **`v1/usbmeta/protocol`** is the base path of the API on which the client obtains the USB class, subclass, and protocol descriptions by providing the class, subclass, and protocol IDs.

#### Path Settings
The **Paths** section of the configuration file specifies directories where various files will be written. Relative paths are prepended with the installation directory.
```json
"Paths": {
    "ReportDir": "report"
}
```
* **`ReportDir`** is where device reports are written. This can be overridden with the `folder` report _option flag_.

#### Logger Settings
The **Loggers** section of the configuration file contains logging options for the system, change, and error log.
```json
"Loggers": {

    "LogDir": "log",
    "Console": false,
    "Syslog": false,
    
    "Logger": {
        "system": {
            "LogFile": "system.log",
            "Console": false,
            "Syslog": false,
            "Prefix": ["date", "time"]
        },
        "change": {
            "LogFile": "change.log",
            "Console": false,
            "Syslog": false,
            "Prefix": ["date", "time"]
        },
        "error": {
            "LogFile": "error.log",
            "Console": true,
            "Syslog": false,
            "Prefix": ["date", "time", "file"]
        }
    }
}
```
* **`LogDir`** is the directory where logs files will be written.
* **`Console`** causes the utility to write events to the console (stdout) in addition to the log file. This overrides the same setting for individual logs, below.
* **`Syslog`** causes the utility to write events to a local or remote syslog daemon using the `Syslog` configuration settings (see _Syslog Settings,_ below).
* **`Logger`** is a collection of logs used by the utility to record events.
    * **`system`** contains settings for the _system log,_ where the utility records significant, non-error events.
    * **`change`** contains settings for the _change log,_ where the utility records changes found during audits. It also reports changes to the server.
    * **`error`** contains settings for the _error log,_ where the utility records critical errors.

Each logger, above, has the following configuration settings:

* **`LogFile`** specifies the filename of the log file.
* **`Console`** specifies whether or not events are written to the console (stdout) in addition to the log file.
* **`Syslog`** causes the utility to write events to a local or remote syslog daemon using the `Syslog` configuration settings (see _Syslog Settings,_ below).
* **`Prefix`** is a comma-separated list of optional attributes that will be prepended to each log entry:
    * **`date`** is the date of the event in _YYYY/MM/DD_ format.
    * **`time`** is the local time of the event in _HH:MM:SS_ format.
    * **`file`** is the name of the file containing the source code that produced the event.

#### Syslog Settings
The **Syslog** section contains parameters for communicating with a local or remote syslog server. Please note that the syslog daemon, if not running on the same host as the utility, must be configured to accept remote syslog client connections.
```json
"Syslog": {
    "Enabled": false,
    "Protocol": "udp",
    "Port": "514",
    "Host": "localhost",
    "Tag": "usbci_cmdbc",
    "Facility": "LOG_LOCAL7",
    "Severity": "LOG_INFO"
}
```
* **`Enabled`** specifies whether or not syslog logging is available to the loggers. If syslog logging is not _enabled,_ the loggers will not write to the configured syslog daemon, even if they're configured to do so.
* **`Protocol`** is the transport-layer protocol used by the syslog daemon (blank for local).
* **`Port`** is the port used by the syslog daemon (blank for local).
* **`Host`** is the hostname or IP address of the syslog daemon (blank for local).
* **`Tag`** is an arbitrary string to prepend to the syslog event.
* **`Facility`** specifies the type of program that is logging the message (see [RFC 5424](https://tools.ietf.org/html/rfc5424)):
    * **`LOG_KERN`** -- kernel messages
    * **`LOG_USER`** -- user-level messages
    * **`LOG_MAIL`** -- mail system
    * **`LOG_DAEMON`** -- system daemons
    * **`LOG_AUTH`** -- security/authorization messages
    * **`LOG_SYSLOG`** -- messages generated internally by syslogd
    * **`LOG_LPR`** -- line printer subsystem
    * **`LOG_NEWS`** -- network news subsystem
    * **`LOG_UUCP`** -- UUCP subsystem
    * **`LOG_CRON`** -- security/authorization messages
    * **`LOG_AUTHPRIV`** -- FTP daemon
    * **`LOG_FTP`** -- scheduling daemon
    * **`LOG_LOCAL0`** -- local use 0
    * **`LOG_LOCAL1`** -- local use 1
    * **`LOG_LOCAL2`** -- local use 2
    * **`LOG_LOCAL3`** -- local use 3
    * **`LOG_LOCAL4`** -- local use 4
    * **`LOG_LOCAL5`** -- local use 5
    * **`LOG_LOCAL6`** -- local use 6
    * **`LOG_LOCAL7`** -- local use 7
* **`Severity`** specifies the severity of the event (see [RFC 5424](https://tools.ietf.org/html/rfc5424)):
    * **`LOG_EMERG`** -- system is unusable
    * **`LOG_ALERT`** -- action must be taken immediately
    * **`LOG_CRIT`** -- critical conditions
    * **`LOG_ERR`** -- error conditions
    * **`LOG_WARNING`** -- warning conditions
    * **`LOG_NOTICE`** -- normal but significant conditions
    * **`LOG_INFO`** -- informational messages
    * **`LOG_DEBUG`** -- debug-level messages


#### Include Settings
The **Include** section specifies device vendors and products to include (_true_) or exclude (_false_) when conducting inventories.
```json
"Include": {
    "VendorID": {
        "0801": true,
        "043d": false,
        "045e": false
    }
    "ProductID": {
        "0acd": {
            "2010": true,
            "2030": true
        },
        "046a": {
            "0001": true
        }
    },
    "Default": false
}
```
* **`VendorID`** specifies which vendors to include or exclude. This setting applies to all of the vendor's products and overrides both the _ProductID_ and _Default_ configuration settings; that is, if a vendor is excluded under _VendorID_, that vendor's products cannot be included under the _ProductID_ or _Default_ sections. Here, all **Magtek** (vendor ID **0801**) devices will be included and **Microsoft** (vendor IDs **043d** and **045e**) devices will be excluded.
* **`ProductID`** specifies individual products to include or exclude. This setting applies to specific _ProductIDs_ under a given _VendorID_ and overrides the _Default_ configuration setting. Here, **IDTech** (vendor ID **0acd**) card readers (product IDs **2010** and **2030** will be included, as will **Cherry** (vendor ID **046a**) keyboards (product ID **0001**). 
* **`Default`** specifies the default behavior for products that are not specifically included or excluded by _Vendor ID_ or _Product ID_. Here the default is to include, which effectively renders previous inclusions redundant; however, specific _VendorID_ and _ProductID_ inclusions ensure that those devices will be inventoried even if the _Default_ setting is changed to 'exclude' (_false_).

#### Format Settings
Default file formats for various use cases:
```json
"Format": {
    "Report": "csv",
    "Default": "json"
}
```
* **`Report`** is the default output format for inventory reports.
* **`Default`** is the default output format for other use cases.

### Command-Line Flags
Client operation is controlled through command-line _flags_. There are seven top-level _action flags_ -- `audit`, `checkin`, `legacy`, `report`, `reset`, `serial`, and `help`.  Some of these require (or offer) additional _option flags_.
* **`-audit`** performs a device configuration change audit.
    * **`-local`** audits against JSON state files stored on the local machine
    * **`-server`**	audits against the last device check-in stored in the database.
    * **`-help`** lists _audit option flags_ and their descriptions.
* **`-checkin`** checks devices in with the server, which stores device information in the database along with the check-in date.
* **`-legacy`** specifies _legacy mode_, which produces the same output to the same filename, `usb_serials.txt`, as the legacy inventory utility. The utility will also operate in legacy mode if the executable is renamed from **cmdbc.exe** to **magtek_inventory.exe**, the name of the legacy inventory utility executable.
* **`-report`** generates device configuration reports.
    * **`-console`** writes report output to the console.
    * **`-folder`** _`<path>`_ writes report output files to _`<path>`_. It defaults to the `report` folder beneath the installation directory.
    * **`-format`** _`<format>`_ specifies which report _`<format>`_ to use:
        * **`csv`** specifies comma-separated value format (default).
        * **`nvp`** specifies name-value pair format.
        * **`xml`** specifies extensible markup language format.
        * **`json`** specifies JavaScript object notation format.
    * **`-help`** lists _report option flags_ and their descriptions.
* **`-reset`** resets the device.
* **`-serial`** performs serial number operations. (By default, **CMDBc** will not configure a serial number on a device that already has one.)
    * **`-copy`** copies the factory serial number (if present) to the active serial number.
    * **`-erase`** erases the current serial number.
    * **`-fetch`** fetches a unique serial number from the server.
    * **`-force`** forces a serial number change, even if the device already has one.
    * **`-set`** _`<value>`_ sets serial number to the specified _`<value>`_.
    * **`-help`** lists _serial option flags_ and their descriptions.
* **`-help`** lists top-level _action flags_ and their descriptions.

### Serial Number Configuration
Configure serial numbers on attached devices with the `serial` _action flag_.

The `set`, `copy`, and `fetch` _option flags_ are mutually-exclusive. You assign a specific serial number string with the `set` _option flag_, copy the immutable factory serial number (if one exists) to the configurable serial number with the `copy` _option flag_, or request a new, unique serial number from the server with the `fetch` _option flag_.

The `copy`, `fetch`, and `set` _option flags_ can each be combined with `erase` and `force`. By default, **CMDBc** ignores serial number changes for devices that already have serial numbers. The `erase` _option flag_ bypasses this by erasing the existing serial number before attempting to assign a new one, effectively removing the constraint. The `force` _option flag_ simply overrides the safeguard feature.

**Examples**:
```sh
cmdbc.exe -serial -fetch -force
```
The preceding command will, for each compatible device, fetch a new serial number from the server and configure the device with it, overriding the safety mechanism that prevents overwriting existing serial numbers.
```sh
cmdbc.exe -serial -erase -fetch
```
The preceding command will, for each compatible device, erase the existing serial number, fetch a new, unique serial number from the server, and configure the device with it.

While the previous two examples would normally produce the same result, a subtle difference is that, if **CMDBc** were unable to obtain a new serial number, `force` would leave existing serial numbers in place whereas `erase` would leave devices without serial numbers.

You can also use the `erase` _option flag_ by itself to erase device serial numbers, although this is an unusual use case.

**Caution**: action and option flags apply to _all attached devices_; if you use the `serial` _action flag_ with the `fetch` _option flag_, **CMDBc** will only configure new serial numbers on compatible devices that don't already have one. If all attached devices already have serial numbers or are not configurable, nothing will happen. However, if you add the `force` flag, it will overwrite the serial number on all compatible devices -- even those that already have a serial number. If you use the `set` and `force` _option flags_ and there is more than one configurable device attached, you will end up having multiple devices with the same serial number.

Refer to the _Database Structure_ section in the documentation for [**CMDBd**](https://github.com/jscherff/cmdbd/blob/master/README.md) for details on device information transferred to the server and tables/columns affected on serial number requests.
 
### Device Registration
Register attached devices with the server using the `checkin` _action flag_. This will create a new object in the device repository.

Refer to the _Database Structure_ section in the documentation for [**CMDBd**](https://github.com/jscherff/cmdbd/blob/master/README.md) for details on device information transferred to the server and tables/columns affected on device registrations.

### Device Audits
Perform a configuration change audit for attached devices using the `audit` _action flag._

You can audit against device state files saved on the local workstation with the `local` _option flag_, or you can audit against device information stored in the database with the `server` _option flag_. The latter is preferred. By default, device state for local audits is stored in JSON files in the `state` subdirectory under the utility installation directory (configurable). Changes detected during an audit are written to the local change log and are also reported to the server.

Audits are only supported on serialized devices.

Refer to the _Database Structure_ section in the documentation for [**CMDBd**](https://github.com/jscherff/cmdbd/blob/master/README.md) for details on device information transferred to the server and tables/columns affected on device audits.

### Device Reports
Generate device reports for attached devices using the `report` _action flag._

Select the report format with the `format` _option flag_. Four formats are currently supported: _comma-separated value_ (CSV), _name-value pairs_ (NVP), _extensible markup language_ (XML), and _JavaScript object notation_ (JSON). 

By default, report files are written to the `report` subdirectory under the utility installation directory (configurable). A separate report file is generated for each device. The report filename is `{bn}-{ba}-{pn}-{vid}-{pid}.{fmt}`, where
* `bn` is a three-digit decimal value representing _bus number_,
* `ba` is a three-digit decimal value representing _bus address_,
* `pn` is a three-digit decimal value representing _port number_,
* `vid` is a four-digit hexadecimal value representing _vendor ID_,
* `pid` is a four-digit hexadecimal value representing _product ID_, and
* `fmt` is the report format (csv, nvp, xml, or json)

Change the report destination folder with the `folder` _option flag_.

Write the report to the console with the `console` _option flag_. If you use the `console` _option flag_ without `folder`, the report will only be written to the console. If you use the `console` _option flag_ after `folder`, the report will be written to the specified folder _and_ to the console. If you use the `console` _option flag_ before `folder`, the report will only be written to the console and `folder` will be ignored.

**Examples**:
```sh
cmdbc.exe -report -format csv
```
The preceding command writes the device reports in CSV format to the 'reports' subdirectory.
```sh
cmdbc.exe -report -format json -console
cmdbc.exe -report -format json -console -folder c:\reports
```
Both of the preceding commands write the device reports in JSON format to the console. The `folder` _option flag_ in the second command is ignored.
```sh
cmdbc.exe -report -format xml -folder c:\reports
cmdbc.exe -report -format xml -folder c:\reports -console
```
Both of the preceding commands write the device reports in XML format to the c:\reports folder. The second command also writes the reports to the console.

### Device Resets
Reset attached devices using the `reset` _action flag_.

Depending on the device, this either does a host-side reset, refreshing the USB device descriptor, or a low-level hardware reset on the device.

### Legacy Reports
Write a legacy device report using the `legacy` _action flag_.

This feature mimics the behavior of previous device inventory utilities for integration backward compatibility. It simply writes the hostname and device serial number in CSV format to a file named `usb_serials.txt` in the utility installation directory, then exits. It filters all but MagTek card readers, and if there is more than one card reader attached, it arbitrarily chooses one.

Renaming the utility from **cmdbd.exe** to **magtek_inventory.exe** forces this behavior without command-line flags.
